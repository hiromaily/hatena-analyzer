package rdb

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb/sqlcgen"
)

// type RDBClient interface {
// 	Close(ctx context.Context) error
// 	Begin(ctx context.Context) error
// 	Commit(ctx context.Context) error
// 	Rollback(ctx context.Context) error
// 	ExecuteSQLFile(ctx context.Context, filepath string) error
// }

//
// Sqlc with PostgreSQL
//

type SqlcPostgresClient struct {
	db        *pgx.Conn
	tx        pgx.Tx
	Queries   *sqlcgen.Queries
	QueriesTx *sqlcgen.Queries
}

func NewSqlcPostgresClient(ctx context.Context, dataSourceName string) (*SqlcPostgresClient, error) {
	// validation
	if dataSourceName == "" {
		return nil, errors.New("dataSourceName is empty")
	}

	db, err := pgx.Connect(ctx, dataSourceName)
	if err != nil {
		return nil, err
	}

	queries := sqlcgen.New(db)

	return &SqlcPostgresClient{
		db:        db,
		tx:        nil,
		Queries:   queries,
		QueriesTx: nil,
	}, nil
}

// Close db connection
func (s *SqlcPostgresClient) Close(ctx context.Context) error {
	if s.db != nil {
		return s.db.Close(ctx)
	}
	return nil
}

func (s *SqlcPostgresClient) Begin(ctx context.Context) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	qtx := s.Queries.WithTx(tx)
	s.tx = tx
	s.QueriesTx = qtx

	return nil
}

func (s *SqlcPostgresClient) Commit(ctx context.Context) error {
	err := s.tx.Commit(ctx)
	if err != nil {
		return err
	}
	s.tx = nil
	s.QueriesTx = nil
	return nil
}

func (s *SqlcPostgresClient) Rollback(ctx context.Context) error {
	err := s.tx.Rollback(ctx)
	if err != nil {
		return err
	}
	s.tx = nil
	s.QueriesTx = nil
	return nil
}

// Execute raw sql
func (s *SqlcPostgresClient) ExecuteSQLFile(ctx context.Context, filepath string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(ctx, string(sqlBytes))
	if err != nil {
		return err
	}
	return nil
}
