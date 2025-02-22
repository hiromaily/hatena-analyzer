package repository

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
)

type URLRepositorier interface {
	Close(ctx context.Context) error
	InsertURLs(ctx context.Context, category entities.CategoryCode, urls []string) error
}

type urlRepository struct {
	logger         logger.Logger
	postgreQueries *rdb.PostgreQueries
}

func NewURLRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
) *urlRepository {
	return &urlRepository{
		logger:         logger,
		postgreQueries: postgreQueries,
	}
}

func (u *urlRepository) Close(ctx context.Context) error {
	return u.postgreQueries.Close(ctx)
}

func (u *urlRepository) InsertURLs(ctx context.Context, category entities.CategoryCode, urls []string) error {
	_, err := u.postgreQueries.InsertURLs(ctx, category, urls)
	return err
}
