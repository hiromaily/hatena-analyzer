package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb/sqlcgen"
)

type UserRepositorier interface {
	Close(ctx context.Context) error
	GetUsers(ctx context.Context) ([]string, error)
	GetUsersByURLS(ctx context.Context, urls []string) ([]string, error)
	UpdateUserBookmarkCount(ctx context.Context, userName string, count int) error
}

type rdbUserRepository struct {
	logger    logger.Logger
	rdbClient *rdb.SqlcPostgresClient
}

func NewRDBUserRepository(
	logger logger.Logger,
	rdbClient *rdb.SqlcPostgresClient,
) *rdbUserRepository {
	return &rdbUserRepository{
		logger:    logger,
		rdbClient: rdbClient,
	}
}

func (r *rdbUserRepository) Close(ctx context.Context) error {
	return r.rdbClient.Close(ctx)
}

func (r *rdbUserRepository) GetUsers(ctx context.Context) ([]string, error) {
	queries, release, err := r.rdbClient.GetQueries(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	return queries.GetUsers(ctx)
}

func (r *rdbUserRepository) GetUsersByURLS(ctx context.Context, urls []string) ([]string, error) {
	queries, release, err := r.rdbClient.GetQueries(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	return queries.GetUsersByURLs(ctx, urls)
}

func (r *rdbUserRepository) UpdateUserBookmarkCount(ctx context.Context, userName string, count int) error {
	param := sqlcgen.UpdateUserBookmarkCountParams{
		BookmarkCount: pgtype.Int4{Int32: int32(count), Valid: true},
		UserName:      userName,
	}
	queries, release, err := r.rdbClient.GetQueries(ctx)
	if err != nil {
		return err
	}
	defer release()
	_, err = queries.UpdateUserBookmarkCount(ctx, param)
	return err
}
