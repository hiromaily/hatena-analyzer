package repository

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/influxdb"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/mongodb"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
)

type FetchBookmarkRepositorier interface {
	Close(ctx context.Context)
	// PostgreSQL
	GetAllURLs(ctx context.Context) ([]entities.URLIDAddress, error)
	GetURLID(ctx context.Context, url string) (int32, error)
	InsertURL(
		ctx context.Context,
		url string,
		categoryCode entities.CategoryCode,
		bmCount, userCount int,
	) (int32, error)
	UpsertURL(
		ctx context.Context,
		url string,
		categoryCode entities.CategoryCode,
		title string,
		bmCount, userCount int,
		privateUserRate float64,
	) (int32, error)
	UpdateURL(
		ctx context.Context,
		urlID int32,
		title string,
		bmCount, userCount int,
		privateUserRate float64,
	) (int64, error)
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
// fetchBookmarkRepository Implementation
//

type fetchBookmarkRepository struct {
	logger          logger.Logger
	postgreQueries  *rdb.PostgreQueries
	influxDBQueries *influxdb.InfluxDBQueries
	mongoDBQueries  *mongodb.MongoDBQueries
}

func NewFetchBookmarkRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
	influxDBQueries *influxdb.InfluxDBQueries,
	mongoDBQueries *mongodb.MongoDBQueries,
) *fetchBookmarkRepository {
	return &fetchBookmarkRepository{
		logger:          logger,
		postgreQueries:  postgreQueries,
		influxDBQueries: influxDBQueries,
		mongoDBQueries:  mongoDBQueries,
	}
}

func (f *fetchBookmarkRepository) Close(ctx context.Context) {
	f.postgreQueries.Close(ctx)
	f.influxDBQueries.Close(ctx)
	f.mongoDBQueries.Close(ctx)
}

// PostgreSQL

func (f *fetchBookmarkRepository) GetAllURLs(ctx context.Context) ([]entities.URLIDAddress, error) {
	return f.postgreQueries.GetAllURLs(ctx)
}

func (f *fetchBookmarkRepository) GetURLID(ctx context.Context, url string) (int32, error) {
	return f.postgreQueries.GetURLID(ctx, url)
}

func (f *fetchBookmarkRepository) InsertURL(
	ctx context.Context,
	url string,
	categoryCode entities.CategoryCode,
	bmCount, userCount int,
) (int32, error) {
	return f.postgreQueries.InsertURL(ctx, url, categoryCode, bmCount, userCount)
}

func (f *fetchBookmarkRepository) UpsertURL(
	ctx context.Context,
	url string,
	categoryCode entities.CategoryCode,
	title string,
	bmCount, userCount int,
	privateUserRate float64,
) (int32, error) {
	return f.postgreQueries.UpsertURL(ctx, url, categoryCode, title, bmCount, userCount, privateUserRate)
}

func (f *fetchBookmarkRepository) UpdateURL(
	ctx context.Context,
	urlID int32,
	title string,
	bmCount, userCount int,
	privateUserRate float64,
) (int64, error) {
	return f.postgreQueries.UpdateURL(ctx, urlID, title, bmCount, userCount, privateUserRate)
}

// func (b *bookmarkRepository) InsertUser(ctx context.Context, userName string) error {
// 	return b.postgreQueries.InsertUser(ctx, userName)
// }

func (f *fetchBookmarkRepository) UpsertUser(ctx context.Context, userName string) (int32, error) {
	return f.postgreQueries.UpsertUser(ctx, userName)
}

func (f *fetchBookmarkRepository) UpsertUserURLs(ctx context.Context, userID, urlID int32) error {
	return f.postgreQueries.UpsertUserURLs(ctx, userID, urlID)
}

// InfluxDB

func (f *fetchBookmarkRepository) ReadEntitySummary(
	ctx context.Context,
	url string,
) (*entities.BookmarkSummary, error) {
	return f.influxDBQueries.ReadEntitySummary(ctx, url)
}

func (f *fetchBookmarkRepository) WriteEntitySummary(
	ctx context.Context,
	url string,
	bookmark *entities.Bookmark,
) error {
	return f.influxDBQueries.WriteEntitySummary(ctx, url, bookmark)
}

// MongoDB

func (f *fetchBookmarkRepository) ReadEntity(ctx context.Context, url string) (*entities.Bookmark, error) {
	return f.mongoDBQueries.ReadEntity(ctx, url)
}

func (f *fetchBookmarkRepository) WriteEntity(
	ctx context.Context,
	url string,
	bookmark *entities.Bookmark,
) error {
	return f.mongoDBQueries.WriteEntity(ctx, url, bookmark)
}
