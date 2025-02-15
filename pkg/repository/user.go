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
	return r.rdbClient.Queries.GetUsers(ctx)
}

func (r *rdbUserRepository) UpdateUserBookmarkCount(ctx context.Context, userName string, count int) error {
	param := sqlcgen.UpdateUserBookmarkCountParams{
		BookmarkCount: pgtype.Int4{Int32: int32(count), Valid: true},
		UserName:      userName,
	}
	_, err := r.rdbClient.Queries.UpdateUserBookmarkCount(ctx, param)
	return err
}
