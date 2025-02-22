package repository

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
)

type UserRepositorier interface {
	Close(ctx context.Context) error
	GetUserNames(ctx context.Context) ([]string, error)
	GetUserNamesByURLS(ctx context.Context, urls []string) ([]string, error)
	UpdateUserBookmarkCount(ctx context.Context, userName string, count int) error
}

type userRepository struct {
	logger         logger.Logger
	postgreQueries *rdb.PostgreQueries
}

func NewUserRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
) *userRepository {
	return &userRepository{
		logger:         logger,
		postgreQueries: postgreQueries,
	}
}

func (u *userRepository) Close(ctx context.Context) error {
	return u.postgreQueries.Close(ctx)
}

func (u *userRepository) GetUserNames(ctx context.Context) ([]string, error) {
	return u.postgreQueries.GetUserNames(ctx)
}

func (u *userRepository) GetUserNamesByURLS(ctx context.Context, urls []string) ([]string, error) {
	return u.postgreQueries.GetUserNamesByURLS(ctx, urls)
}

func (u *userRepository) UpdateUserBookmarkCount(ctx context.Context, userName string, count int) error {
	return u.postgreQueries.UpdateUserBookmarkCount(ctx, userName, count)
}
