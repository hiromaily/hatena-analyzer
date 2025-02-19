package influxdb

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
)

type InfluxDBQueries struct {
	logger   logger.Logger
	dbClient influxdb2.Client // Client interface
	org      string
	bucket   string
}

func NewInfluxDBQueries(
	logger logger.Logger,
	dbClient influxdb2.Client,
	org, bucket string,
) *InfluxDBQueries {
	return &InfluxDBQueries{
		logger:   logger,
		dbClient: dbClient,
		org:      org,
		bucket:   bucket,
	}
}

func (i *InfluxDBQueries) Close(_ context.Context) {
	i.dbClient.Close()
}

func (i *InfluxDBQueries) ReadEntitySummary(
	ctx context.Context,
	url string,
) (*entities.BookmarkSummary, error) {
	// query
	queryAPI := i.dbClient.QueryAPI(i.org)
	query := fmt.Sprintf(`
	from(bucket: "%s") 
	|> range(start: -1d) 
	|> filter(fn: (r) => r._measurement == "%s")
	|> filter(fn: (r) => r._field == "count" or r._field == "user_num" or r._field == "deleted_user_num")
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
	var (
		latestCount          int
		latestUserNum        int
		latestDeletedUserNum int
		timeStamp            time.Time
	)

	for result.Next() {
		record := result.Record()
		timeStamp = record.Time()

		switch record.Field() {
		case "count":
			countValue, ok := record.Value().(int64)
			if !ok {
				i.logger.Error("expecting count to be int64")
			}
			latestCount = int(countValue)
		case "user_num":
			userNumValue, ok := record.Value().(int64)
			if !ok {
				i.logger.Error("expecting user_num to be int64")
			}
			latestUserNum = int(userNumValue)
		case "deleted_user_num":
			userNumValue, ok := record.Value().(int64)
			if !ok {
				i.logger.Error("expecting deleted user_num to be int64")
			}
			latestDeletedUserNum = int(userNumValue)
		default:
			i.logger.Error("unexpected field", "field", record.Field())
		}
	}
	if result.Err() != nil {
		i.logger.Error("failed to retrieve data", "error", result.Err())
		return nil, result.Err()
	}

	i.logger.Debug("latest point",
		"time", timeStamp,
		"count", latestCount,
		"user_num", latestUserNum,
		"deleted_user_num", latestDeletedUserNum,
	)

	bookmarkSummary := &entities.BookmarkSummary{
		Count:            latestCount,
		UserCount:        latestUserNum,
		DeletedUserCount: latestDeletedUserNum,
		Timestamp:        timeStamp,
	}

	return bookmarkSummary, nil
}

func (i *InfluxDBQueries) ReadEntitySummaries(
	ctx context.Context,
	url string,
) ([]*entities.BookmarkSummary, error) {
	// query
	queryAPI := i.dbClient.QueryAPI(i.org)
	query := fmt.Sprintf(`
		from(bucket: "%s") 
		|> range(start: -1d) 
		|> filter(fn: (r) => r._measurement == "%s")
		|> filter(fn: (r) => r._field == "count" or r._field == "user_num" or r._field == "deleted_user_num")
		|> sort(columns: ["_time"], desc: true)
		`, i.bucket, url)

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		// Debug: what happened when data is not found
		i.logger.Error("failed to call influxDB queryAPI.Query()", "url", url, "error", err)
		return nil, err
	}

	// initialize slice for summaries
	var summaries []*entities.BookmarkSummary

	var recordsMap = make(map[time.Time]*entities.BookmarkSummary)

	for result.Next() {
		record := result.Record()
		timeStamp := record.Time()

		// Get title tag
		title, ok := record.ValueByKey("title").(string)
		if !ok {
			i.logger.Error("expecting title to be string")
			title = "unknown" // fallback value
		}

		// For each record, retrieve corresponding BookmarkSummary or create a new one if it doesn't exist
		bookmarkSummary, ok := recordsMap[timeStamp]
		if !ok {
			bookmarkSummary = &entities.BookmarkSummary{
				Timestamp: timeStamp,
				Title:     title,
			}
			recordsMap[timeStamp] = bookmarkSummary
		}

		switch record.Field() {
		case "count":
			countValue, ok := record.Value().(int64)
			if !ok {
				i.logger.Error("expecting count to be int64")
			}
			bookmarkSummary.Count = int(countValue)
		case "user_num":
			userNumValue, ok := record.Value().(int64)
			if !ok {
				i.logger.Error("expecting user_num to be int64")
			}
			bookmarkSummary.UserCount = int(userNumValue)
		case "deleted_user_num":
			userNumValue, ok := record.Value().(int64)
			if !ok {
				i.logger.Error("expecting deleted user_num to be int64")
			}
			bookmarkSummary.DeletedUserCount = int(userNumValue)
		default:
			i.logger.Error("unexpected field", "field", record.Field())
		}
	}
	if result.Err() != nil {
		i.logger.Error("failed to retrieve data", "error", result.Err())
		return nil, result.Err()
	}

	// Convert map to slice
	for _, value := range recordsMap {
		summaries = append(summaries, value)
	}

	// sort summaries by time
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Timestamp.Before(summaries[j].Timestamp)
	})

	return summaries, nil
}

func (i *InfluxDBQueries) WriteEntitySummary(
	ctx context.Context,
	url string,
	bookmark *entities.Bookmark,
) error {
	if bookmark == nil {
		return errors.New("bookmark is nil")
	}

	userNum := len(bookmark.Users)
	deletedUserNum := bookmark.CountDeletedUser()

	writeAPI := i.dbClient.WriteAPIBlocking(i.org, i.bucket)

	i.logger.Debug(
		"data will be stored",
		"title", bookmark.Title,
		"count", bookmark.Count,
		"user_num", userNum,
		"deleted_user_num", deletedUserNum,
	)

	point := influxdb2.NewPointWithMeasurement(url).
		AddTag("title", bookmark.Title).
		AddField("count", bookmark.Count).
		AddField("user_num", userNum).
		AddField("deleted_user_num", deletedUserNum).
		SetTime(time.Now())

	return writeAPI.WritePoint(ctx, point)
}
