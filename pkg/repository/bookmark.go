package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookmarkRepositorier interface {
	Close(ctx context.Context)
	ReadEntitySummary(ctx context.Context, url string) (*entities.BookmarkSummary, error)
	WriteEntitySummary(ctx context.Context, url string, bookmark *entities.Bookmark) error
	ReadEntity(ctx context.Context, url string) (*entities.Bookmark, error)
	WriteEntity(ctx context.Context, url string, bookmark *entities.Bookmark) error
}

type bookmarkRepository struct {
	logger               logger.Logger
	influxDBBookmarkRepo *influxDBBookmarkRepository
	mongoDBBookmarkRepo  *mongoDBBookmarkRepository
}

func NewBookmarkRepository(
	logger logger.Logger,
	influxDBBookmarkRepo *influxDBBookmarkRepository,
	mongoDBBookmarkRepo *mongoDBBookmarkRepository,
) *bookmarkRepository {
	return &bookmarkRepository{
		logger:               logger,
		influxDBBookmarkRepo: influxDBBookmarkRepo,
		mongoDBBookmarkRepo:  mongoDBBookmarkRepo,
	}
}

func (b *bookmarkRepository) Close(ctx context.Context) {
	b.influxDBBookmarkRepo.Close(ctx)
	b.mongoDBBookmarkRepo.Close(ctx)
}

func (b *bookmarkRepository) ReadEntitySummary(
	ctx context.Context,
	url string,
) (*entities.BookmarkSummary, error) {
	return b.influxDBBookmarkRepo.ReadEntitySummary(ctx, url)
}

func (b *bookmarkRepository) WriteEntitySummary(
	ctx context.Context,
	url string,
	bookmark *entities.Bookmark,
) error {
	return b.influxDBBookmarkRepo.WriteEntitySummary(ctx, url, bookmark)
}

func (b *bookmarkRepository) ReadEntity(ctx context.Context, url string) (*entities.Bookmark, error) {
	return b.mongoDBBookmarkRepo.ReadEntity(ctx, url)
}

func (b *bookmarkRepository) WriteEntity(ctx context.Context, url string, bookmark *entities.Bookmark) error {
	return b.mongoDBBookmarkRepo.WriteEntity(ctx, url, bookmark)
}

//
// InfluxDBBookmarkRepository Implementation
//

type influxDBBookmarkRepository struct {
	logger   logger.Logger
	dbClient influxdb2.Client // Client interface
	org      string
	bucket   string
}

func NewInfluxDBBookmarkRepository(
	logger logger.Logger,
	dbClient influxdb2.Client,
	org, bucket string,
) *influxDBBookmarkRepository {
	return &influxDBBookmarkRepository{
		logger:   logger,
		dbClient: dbClient,
		org:      org,
		bucket:   bucket,
	}
}

func (i *influxDBBookmarkRepository) Close(_ context.Context) {
	i.dbClient.Close()
}

func (i *influxDBBookmarkRepository) ReadEntitySummary(
	ctx context.Context,
	url string,
) (*entities.BookmarkSummary, error) {
	// query
	queryAPI := i.dbClient.QueryAPI(i.org)
	query := fmt.Sprintf(`
	from(bucket: "%s") 
	|> range(start: -1d) 
	|> filter(fn: (r) => r._measurement == "%s")
	|> filter(fn: (r) => r._field == "count" or r._field == "user_num")
	|> sort(columns: ["_time"], desc: true)
	|> limit(n: 1)
	`, i.bucket, url)

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		// Debug: what happened when data is not found
		i.logger.Error("failed to call influxDB queryAPI.Query()", "url", url, "error", err)
		return nil, err
	}

	// retrieve data
	var latestCount int
	var latestUserNum int
	var timeStamp time.Time

	for result.Next() {
		record := result.Record()
		timeStamp = record.Time()

		if record.Field() == "count" {
			countValue, ok := record.Value().(int64)
			if !ok {
				i.logger.Error("expecting count to be int64")
			}
			latestCount = int(countValue)
		} else if record.Field() == "user_num" {
			userNumValue, ok := record.Value().(int64)
			if !ok {
				i.logger.Error("expecting count to be int64")
			}
			latestUserNum = int(userNumValue)
		} else {
			i.logger.Error("unexpected field", "field", record.Field())
		}
	}
	if result.Err() != nil {
		i.logger.Error("failed to retrieve data", "error", result.Err())
		return nil, result.Err()
	}

	i.logger.Debug("latest point", "time", timeStamp, "count", latestCount, "user_num", latestUserNum)

	bookmarkSummary := &entities.BookmarkSummary{
		Count:     latestCount,
		UserCount: latestUserNum,
		Timestamp: timeStamp,
	}

	return bookmarkSummary, nil
}

func (i *influxDBBookmarkRepository) WriteEntitySummary(
	ctx context.Context,
	url string,
	bookmark *entities.Bookmark,
) error {
	if bookmark == nil {
		return errors.New("bookmark is nil")
	}

	writeAPI := i.dbClient.WriteAPIBlocking(i.org, i.bucket)

	i.logger.Debug(
		"data will be stored",
		"title",
		bookmark.Title,
		"count",
		bookmark.Count,
		"user_num",
		len(bookmark.Users),
	)

	point := influxdb2.NewPointWithMeasurement(url).
		AddTag("title", bookmark.Title).
		AddField("count", bookmark.Count).
		AddField("user_num", len(bookmark.Users)).
		SetTime(time.Now())

	if err := writeAPI.WritePoint(ctx, point); err != nil {
		return err
	}

	return nil
}

//
// MongoDBBookmarkRepository Implementation
//

type mongoDBBookmarkRepository struct {
	logger logger.Logger
	// Mongodb
	mongoClient    *mongo.Client
	mongoDB        *mongo.Database
	collectionName string
}

func NewMongoDBBookmarkRepository(
	logger logger.Logger,
	mongoClient *mongo.Client,
	dbName, collectionName string,
) *mongoDBBookmarkRepository {
	db := mongoClient.Database(dbName)

	return &mongoDBBookmarkRepository{
		logger:         logger,
		mongoClient:    mongoClient,
		mongoDB:        db,
		collectionName: collectionName,
	}
}

func (m *mongoDBBookmarkRepository) Close(ctx context.Context) {
	//nolint:errcheck
	m.mongoClient.Disconnect(ctx)
}

type URLBookmarkDocument struct {
	URL string            `bson:"_id"`
	Doc entities.Bookmark `bson:"doc"`
}

func (m *mongoDBBookmarkRepository) ReadEntity(
	ctx context.Context,
	url string,
) (*entities.Bookmark, error) {
	// create collection
	collection := m.mongoDB.Collection(m.collectionName)

	// find
	var doc URLBookmarkDocument
	err := collection.FindOne(context.TODO(), bson.M{"_id": url}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// not found
			return nil, nil
		}
		return nil, err
	}

	return &doc.Doc, nil
}

func (m *mongoDBBookmarkRepository) WriteEntity(
	ctx context.Context,
	url string,
	bookmark *entities.Bookmark,
) error {
	if bookmark == nil {
		return errors.New("bookmark is nil")
	}

	// create collection
	collection := m.mongoDB.Collection(m.collectionName)

	// insert
	// _, err := collection.InsertOne(ctx, URLBookmarkDocument{
	// 	URL: url,
	// 	Doc: *bookmark,
	// })

	// upsert
	isUpsert := true
	_, err := collection.UpdateOne(ctx, bson.M{"_id": url}, bson.D{
		{Key: "$set", Value: URLBookmarkDocument{
			URL: url,
			Doc: *bookmark,
		}},
	}, &options.UpdateOptions{Upsert: &isUpsert})
	if err != nil {
		return err
	}

	return nil
}
