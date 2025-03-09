package repository

import (
	"context"

	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/storage/influxdb"
	"github.com/hiromaily/hatena-analyzer/pkg/storage/mongodb"
	"github.com/hiromaily/hatena-analyzer/pkg/storage/rdb"
)

type CloserRepositorier interface {
	Close(ctx context.Context)
}

type closerRepository struct {
	logger          logger.Logger
	postgreQueries  *rdb.PostgreQueries
	influxDBQueries *influxdb.InfluxDBQueries
	mongoDBQueries  *mongodb.MongoDBQueries
}

func NewCloserRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
	influxDBQueries *influxdb.InfluxDBQueries,
	mongoDBQueries *mongodb.MongoDBQueries,
) *closerRepository {
	return &closerRepository{
		logger:          logger,
		postgreQueries:  postgreQueries,
		influxDBQueries: influxDBQueries,
		mongoDBQueries:  mongoDBQueries,
	}
}

func (c *closerRepository) Close(ctx context.Context) {
	c.postgreQueries.Close(ctx)
	c.influxDBQueries.Close(ctx)
	c.mongoDBQueries.Close(ctx)
}
