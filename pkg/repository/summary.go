package repository

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/influxdb"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
)

type SummaryRepositorier interface {
	Close(ctx context.Context)
	ReadEntitySummaries(ctx context.Context, url string) ([]*entities.BookmarkSummary, error)
	GetAllURLs(ctx context.Context) ([]entities.RDBURL, error)
	GetUsersByURL(ctx context.Context, url string) ([]entities.RDBUser, error)
}

//
// summaryRepository Implementation
//

type summaryRepository struct {
	logger          logger.Logger
	postgreQueries  *rdb.PostgreQueries
	influxDBQueries *influxdb.InfluxDBQueries
}

func NewSummaryRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
	influxDBQueries *influxdb.InfluxDBQueries,
) *summaryRepository {
	return &summaryRepository{
		logger:          logger,
		postgreQueries:  postgreQueries,
		influxDBQueries: influxDBQueries,
	}
}

func (s *summaryRepository) Close(ctx context.Context) {
	s.postgreQueries.Close(ctx)
	s.influxDBQueries.Close(ctx)
}

// InfluxDB

func (s *summaryRepository) ReadEntitySummaries(
	ctx context.Context,
	url string,
) ([]*entities.BookmarkSummary, error) {
	return s.influxDBQueries.ReadEntitySummaries(ctx, url)
}

// PostgreSQL

func (s *summaryRepository) GetAllURLs(ctx context.Context) ([]entities.RDBURL, error) {
	return s.postgreQueries.GetAllURLs(ctx)
}

func (s *summaryRepository) GetUsersByURL(ctx context.Context, url string) ([]entities.RDBUser, error) {
	return s.postgreQueries.GetUsersByURL(ctx, url)
}
