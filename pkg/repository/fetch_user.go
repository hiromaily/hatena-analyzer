package repository

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
)

type FetchUserRepositorier interface {
	Close(ctx context.Context) error
	GetUserNames(ctx context.Context) ([]string, error)
	GetUserNamesByURLS(ctx context.Context, urls []string) ([]string, error)
	UpdateUserBookmarkCount(ctx context.Context, userName string, count int) error
}

type fetchUserRepository struct {
	logger         logger.Logger
	postgreQueries *rdb.PostgreQueries
}

func NewFetchUserRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
) *fetchUserRepository {
	return &fetchUserRepository{
		logger:         logger,
		postgreQueries: postgreQueries,
	}
}

func (f *fetchUserRepository) Close(ctx context.Context) error {
	return f.postgreQueries.Close(ctx)
}

func (f *fetchUserRepository) GetUserNames(ctx context.Context) ([]string, error) {
	return f.postgreQueries.GetUserNames(ctx)
}

func (f *fetchUserRepository) GetUserNamesByURLS(ctx context.Context, urls []string) ([]string, error) {
	return f.postgreQueries.GetUserNamesByURLS(ctx, urls)
}

func (f *fetchUserRepository) UpdateUserBookmarkCount(ctx context.Context, userName string, count int) error {
	return f.postgreQueries.UpdateUserBookmarkCount(ctx, userName, count)
}
