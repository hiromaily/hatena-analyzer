package repository

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
)

type BookmarkDetailsRepositorier interface {
	Close(ctx context.Context)
	GetAllURLs(ctx context.Context) ([]entities.RDBURL, error)
	GetUsersByURL(ctx context.Context, url string) ([]entities.RDBUser, error)
}

//
// bookmarkDetailsRepository Implementation
//

type bookmarkDetailsRepository struct {
	logger         logger.Logger
	postgreQueries *rdb.PostgreQueries
}

func NewBookmarkDetailsRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
) *bookmarkDetailsRepository {
	return &bookmarkDetailsRepository{
		logger:         logger,
		postgreQueries: postgreQueries,
	}
}

func (b *bookmarkDetailsRepository) Close(ctx context.Context) {
	b.postgreQueries.Close(ctx)
}

// PostgreSQL

func (b *bookmarkDetailsRepository) GetAllURLs(ctx context.Context) ([]entities.RDBURL, error) {
	return b.postgreQueries.GetAllURLs(ctx)
}

func (b *bookmarkDetailsRepository) GetUsersByURL(
	ctx context.Context,
	url string,
) ([]entities.RDBUser, error) {
	return b.postgreQueries.GetUsersByURL(ctx, url)
}
