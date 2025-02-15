package rdb

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb/sqlcgen"
)

//
// Sqlc with PostgreSQL
//

func IsNoRows(err error) bool {
	// "no rows in result set"
	return errors.Is(err, pgx.ErrNoRows)
}

type SqlcPostgresClient struct {
	pool *pgxpool.Pool // pool for concurrent connection
	tx   pgx.Tx        // created when Begin() is called
}

func NewSqlcPostgresClient(
	ctx context.Context,
	dataSourceName string,
	maxConnection int32,
) (*SqlcPostgresClient, error) {
	// validation
	if dataSourceName == "" {
		return nil, errors.New("dataSourceName is empty")
	}
	if maxConnection == 0 {
		return nil, errors.New("maxConnection doesn't allow 0")
	}

	// use Pool
	config, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		return nil, err
	}
	// set config
	config.MaxConns = maxConnection

	// create pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &SqlcPostgresClient{
		pool: pool,
		tx:   nil,
	}, nil
}

func (s *SqlcPostgresClient) GetQueries(ctx context.Context) (*sqlcgen.Queries, func(), error) {
	queries, release, err := func() (*sqlcgen.Queries, func(), error) {
		// use pool
		conn, err := s.pool.Acquire(ctx)
		if err != nil {
			return nil, nil, err
		}
		return sqlcgen.New(conn), conn.Release, nil
	}()
	if err != nil {
		return nil, nil, err
	}
	if s.tx != nil {
		return queries.WithTx(s.tx), release, nil
	}
	return queries, release, nil
}

// Close db connection
func (s *SqlcPostgresClient) Close(_ context.Context) error {
	if s.pool != nil {
		s.pool.Close()
	}
	return nil
}

func (s *SqlcPostgresClient) Begin(ctx context.Context) error {
	tx, err := s.pool.Begin(ctx) // FIXME: maybe wrong
	if err != nil {
		return err
	}
	s.tx = tx
	return nil
}

func (s *SqlcPostgresClient) Commit(ctx context.Context) error {
	err := s.tx.Commit(ctx)
	if err != nil {
		return err
	}
	s.tx = nil
	return nil
}

func (s *SqlcPostgresClient) Rollback(ctx context.Context) error {
	err := s.tx.Rollback(ctx)
	if err != nil {
		return err
	}
	s.tx = nil
	return nil
}

// Execute raw sql
func (s *SqlcPostgresClient) ExecuteSQLFile(ctx context.Context, filepath string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = s.pool.Exec(ctx, string(sqlBytes))
	if err != nil {
		return err
	}
	return nil
}
