package repository

import (
	"context"

	"github.com/hiromaily/hatena-analyzer/pkg/entities"
	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/storage/rdb"
)

type FetchURLRepositorier interface {
	Close(ctx context.Context) error
	// InsertURLs(ctx context.Context, category entities.CategoryCode, urls []string) error
	CallBulkInsertURLs(
		ctx context.Context,
		urls []string,
		categories []entities.CategoryCode,
		isAlls []bool,
	) error
}

type fetchURLRepository struct {
	logger         logger.Logger
	postgreQueries *rdb.PostgreQueries
}

func NewFetchURLRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
) *fetchURLRepository {
	return &fetchURLRepository{
		logger:         logger,
		postgreQueries: postgreQueries,
	}
}

func (f *fetchURLRepository) Close(ctx context.Context) error {
	return f.postgreQueries.Close(ctx)
}

// func (f *fetchURLRepository) InsertURLs(
// 	ctx context.Context,
// 	category entities.CategoryCode,
// 	urls []string,
// ) error {
// 	_, err := f.postgreQueries.InsertURLs(ctx, category, urls)
// 	return err
// }

func (f *fetchURLRepository) CallBulkInsertURLs(
	ctx context.Context,
	urls []string,
	categories []entities.CategoryCode,
	isAlls []bool,
) error {
	return f.postgreQueries.CallBulkInsertURLs(ctx, urls, categories, isAlls)
}
