package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
)

type MongoDBQueries struct {
	logger logger.Logger
	// Mongodb
	mongoClient    *mongo.Client
	mongoDB        *mongo.Database
	collectionName string
}

func NewMongoDBQueries(
	logger logger.Logger,
	mongoClient *mongo.Client,
	dbName, collectionName string,
) *MongoDBQueries {
	db := mongoClient.Database(dbName)

	return &MongoDBQueries{
		logger:         logger,
		mongoClient:    mongoClient,
		mongoDB:        db,
		collectionName: collectionName,
	}
}

func (m *MongoDBQueries) Close(ctx context.Context) {
	//nolint:errcheck
	m.mongoClient.Disconnect(ctx)
}

type URLBookmarkDocument struct {
	URL string            `bson:"_id"`
	Doc entities.Bookmark `bson:"doc"`
}

func (m *MongoDBQueries) ReadEntity(
	ctx context.Context,
	url string,
) (*entities.Bookmark, error) {
	// create collection
	collection := m.mongoDB.Collection(m.collectionName)

	// find
	var doc URLBookmarkDocument
	err := collection.FindOne(ctx, bson.M{"_id": url}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// not found
			return nil, nil
		}
		return nil, err
	}

	return &doc.Doc, nil
}

func (m *MongoDBQueries) WriteEntity(
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
