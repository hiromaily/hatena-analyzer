package repository

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/influxdb"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/mongodb"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
)

type BookmarkRepositorier interface {
	Close(ctx context.Context)
	// PostgreSQL
	GetAllURLs(ctx context.Context) ([]entities.RDBURL, error)
	GetURLID(ctx context.Context, url string) (int32, error)
	InsertURL(ctx context.Context, url string, categoryCode entities.CategoryCode) (int32, error)
	// InsertUser(ctx context.Context, userName string) error
	UpsertUser(ctx context.Context, userName string) (int32, error)
	UpsertUserURLs(ctx context.Context, userID, urlID int32) error
	// InfluxDB
	ReadEntitySummary(ctx context.Context, url string) (*entities.BookmarkSummary, error)
	WriteEntitySummary(ctx context.Context, url string, bookmark *entities.Bookmark) error
	// MongoDB
	ReadEntity(ctx context.Context, url string) (*entities.Bookmark, error)
	WriteEntity(ctx context.Context, url string, bookmark *entities.Bookmark) error
}

//
// bookmarkRepository Implementation
//

type bookmarkRepository struct {
	logger          logger.Logger
	postgreQueries  *rdb.PostgreQueries
	influxDBQueries *influxdb.InfluxDBQueries
	mongoDBQueries  *mongodb.MongoDBQueries
}

func NewBookmarkRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
	influxDBQueries *influxdb.InfluxDBQueries,
	mongoDBQueries *mongodb.MongoDBQueries,
) *bookmarkRepository {
	return &bookmarkRepository{
		logger:          logger,
		postgreQueries:  postgreQueries,
		influxDBQueries: influxDBQueries,
		mongoDBQueries:  mongoDBQueries,
	}
}

func (b *bookmarkRepository) Close(ctx context.Context) {
	b.postgreQueries.Close(ctx)
	b.influxDBQueries.Close(ctx)
	b.mongoDBQueries.Close(ctx)
}

// PostgreSQL

func (b *bookmarkRepository) GetAllURLs(ctx context.Context) ([]entities.RDBURL, error) {
	return b.postgreQueries.GetAllURLs(ctx)
}

func (b *bookmarkRepository) GetURLID(ctx context.Context, url string) (int32, error) {
	return b.postgreQueries.GetURLID(ctx, url)
}

func (b *bookmarkRepository) InsertURL(
	ctx context.Context,
	url string,
	categoryCode entities.CategoryCode,
) (int32, error) {
	return b.postgreQueries.InsertURL(ctx, url, categoryCode)
}

// func (b *bookmarkRepository) InsertUser(ctx context.Context, userName string) error {
// 	return b.postgreQueries.InsertUser(ctx, userName)
// }

func (b *bookmarkRepository) UpsertUser(ctx context.Context, userName string) (int32, error) {
	return b.postgreQueries.UpsertUser(ctx, userName)
}

func (b *bookmarkRepository) UpsertUserURLs(ctx context.Context, userID, urlID int32) error {
	return b.postgreQueries.UpsertUserURLs(ctx, userID, urlID)
}

// InfluxDB

func (b *bookmarkRepository) ReadEntitySummary(
	ctx context.Context,
	url string,
) (*entities.BookmarkSummary, error) {
	return b.influxDBQueries.ReadEntitySummary(ctx, url)
}

func (b *bookmarkRepository) WriteEntitySummary(
	ctx context.Context,
	url string,
	bookmark *entities.Bookmark,
) error {
	return b.influxDBQueries.WriteEntitySummary(ctx, url, bookmark)
}

// MongoDB

func (b *bookmarkRepository) ReadEntity(ctx context.Context, url string) (*entities.Bookmark, error) {
	return b.mongoDBQueries.ReadEntity(ctx, url)
}

func (b *bookmarkRepository) WriteEntity(ctx context.Context, url string, bookmark *entities.Bookmark) error {
	return b.mongoDBQueries.WriteEntity(ctx, url, bookmark)
}
