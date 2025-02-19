package repository

import (
	"context"
	"fmt"
	"sort"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/hiromaily/hatena-fake-detector/pkg/adapter"
	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
)

type SummaryRepositorier interface {
	Close(ctx context.Context)
	ReadEntitySummaries(ctx context.Context, url string) ([]*entities.BookmarkSummary, error)
	GetUsersByURL(ctx context.Context, url string) ([]entities.RDBUser, error)
}

//
// summaryRepository Implementation
//

type summaryRepository struct {
	logger              logger.Logger
	rdbSummaryRepo      *rdbSummaryRepository
	influxDBSummaryRepo *influxDBSummaryRepository
}

func NewSummaryRepository(
	logger logger.Logger,
	rdbSummaryRepo *rdbSummaryRepository,
	influxDBSummaryRepo *influxDBSummaryRepository,
) *summaryRepository {
	return &summaryRepository{
		logger:              logger,
		rdbSummaryRepo:      rdbSummaryRepo,
		influxDBSummaryRepo: influxDBSummaryRepo,
	}
}

func (s *summaryRepository) Close(ctx context.Context) {
	s.rdbSummaryRepo.Close(ctx)
	s.influxDBSummaryRepo.Close(ctx)
}

// InfluxDB

func (s *summaryRepository) ReadEntitySummaries(
	ctx context.Context,
	url string,
) ([]*entities.BookmarkSummary, error) {
	return s.influxDBSummaryRepo.ReadEntitySummaries(ctx, url)
}

// PostgreSQL

func (s *summaryRepository) GetUsersByURL(ctx context.Context, url string) ([]entities.RDBUser, error) {
	return s.rdbSummaryRepo.GetUsersByURL(ctx, url)
}

//
// InfluxDBSummaryRepository Implementation
//

type influxDBSummaryRepository struct {
	logger   logger.Logger
	dbClient influxdb2.Client // Client interface
	org      string
	bucket   string
}

func NewInfluxDBSummaryRepository(
	logger logger.Logger,
	dbClient influxdb2.Client,
	org, bucket string,
) *influxDBSummaryRepository {
	return &influxDBSummaryRepository{
		logger:   logger,
		dbClient: dbClient,
		org:      org,
		bucket:   bucket,
	}
}

func (i *influxDBSummaryRepository) Close(_ context.Context) {
	i.dbClient.Close()
}

func (i *influxDBSummaryRepository) ReadEntitySummaries(
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

//
// PostgresSummaryRepository Implementation
//

type rdbSummaryRepository struct {
	logger    logger.Logger
	rdbClient *rdb.SqlcPostgresClient
}

func NewRDBSummaryRepository(
	logger logger.Logger,
	rdbClient *rdb.SqlcPostgresClient,
) *rdbSummaryRepository {
	return &rdbSummaryRepository{
		logger:    logger,
		rdbClient: rdbClient,
	}
}

func (r *rdbSummaryRepository) Close(ctx context.Context) error {
	return r.rdbClient.Close(ctx)
}

func (r *rdbSummaryRepository) GetUsersByURL(ctx context.Context, url string) ([]entities.RDBUser, error) {
	queries, release, err := r.rdbClient.GetQueries(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	users, err := queries.GetUsersByURL(ctx, url)
	if err != nil {
		return nil, err
	}
	// convert to entity models
	return adapter.UserDBToEntityModel(users), nil
}
