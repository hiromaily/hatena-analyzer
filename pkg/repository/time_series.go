package repository

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/influxdb"
)

type TimeSeriesRepositorier interface {
	Close(ctx context.Context)
	ReadEntitySummaries(ctx context.Context, url string) ([]*entities.BookmarkSummary, error)
}

//
// timeSeriesRepository Implementation
//

type timeSeriesRepository struct {
	logger          logger.Logger
	influxDBQueries *influxdb.InfluxDBQueries
}

func NewTimeSeriesRepository(
	logger logger.Logger,
	influxDBQueries *influxdb.InfluxDBQueries,
) *timeSeriesRepository {
	return &timeSeriesRepository{
		logger:          logger,
		influxDBQueries: influxDBQueries,
	}
}

func (s *timeSeriesRepository) Close(ctx context.Context) {
	s.influxDBQueries.Close(ctx)
}

// InfluxDB

func (s *timeSeriesRepository) ReadEntitySummaries(
	ctx context.Context,
	url string,
) ([]*entities.BookmarkSummary, error) {
	return s.influxDBQueries.ReadEntitySummaries(ctx, url)
}
