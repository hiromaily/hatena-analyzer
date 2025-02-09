package repository

import (
	"context"
	"fmt"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type BookmarkRepositorier interface {
	Close()
	ReadEntity(ctx context.Context, url string) (*entities.Bookmark, error)
	WriteEntity(ctx context.Context, url string, bookmark *entities.Bookmark) error
}

//
// InfluxDBBookmarkRepository
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

func (i *influxDBBookmarkRepository) Close() {
	i.dbClient.Close()
}

func (i *influxDBBookmarkRepository) ReadEntity(ctx context.Context, url string) (*entities.Bookmark, error) {
	var bookmark entities.Bookmark
	bookmark.Users = make(map[string]entities.User)

	queryAPI := i.dbClient.QueryAPI(i.org)
	query := fmt.Sprintf(`
	from(bucket: "%s")
	  |> range(start: 0)
	  |> filter(fn: (r) => r._measurement == "user" and r.url == "%s")
	  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	`, i.bucket, url)

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return &bookmark, err
	}

	for result.Next() {
		record := result.Record()
		// Debug: 詳細なログ出力を追加
		//fmt.Printf("Record: %+v\n", record)

		// 必要なキーが存在するかをチェックし、適切な型にアサートする
		userName, ok := record.ValueByKey("name").(string)
		if !ok {
			fmt.Println("name field missing or not a string")
			continue
		}

		isDeleted, _ := record.ValueByKey("is_deleted").(bool)
		isCommented, _ := record.ValueByKey("is_commented").(bool)

		bookmark.Users[userName] = entities.User{
			Name:        userName,
			IsDeleted:   isDeleted,
			IsCommented: isCommented,
		}
	}
	if result.Err() != nil {
		return &bookmark, result.Err()
	}

	return &bookmark, nil
}

func (i *influxDBBookmarkRepository) WriteEntity(
	ctx context.Context,
	url string,
	bookmark *entities.Bookmark,
) error {
	writeAPI := i.dbClient.WriteAPIBlocking(i.org, i.bucket)

	tags := map[string]string{"url": url}
	fields := map[string]interface{}{
		"title": bookmark.Title,
		"count": bookmark.Count,
	}

	// Bookmarkデータポイントの作成
	point := write.NewPoint("bookmark", tags, fields, bookmark.Timestamp) // with timestamp

	err := writeAPI.WritePoint(ctx, point)
	if err != nil {
		return err
	}

	// Userデータポイントの作成
	for _, user := range bookmark.Users {
		userFields := map[string]interface{}{
			"name":         user.Name,
			"is_commented": user.IsCommented,
			"is_deleted":   user.IsDeleted,
		}
		userPoint := write.NewPoint("user", tags, userFields, bookmark.Timestamp)
		err = writeAPI.WritePoint(context.Background(), userPoint)
		if err != nil {
			return err
		}
	}

	return nil
}
