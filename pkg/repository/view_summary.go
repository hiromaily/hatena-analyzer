package repository

import (
	"context"

	"github.com/hiromaily/hatena-analyzer/pkg/entities"
	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/storage/rdb"
)

type SummaryRepositorier interface {
	Close(ctx context.Context)
	GetAllURLs(ctx context.Context) ([]entities.URL, error)
	GetURLsByURLAddresses(ctx context.Context, urls []string) ([]entities.URL, error)
	// GetUsersByURL(ctx context.Context, url string) ([]entities.RDBUser, error)
	GetAveragePrivateUserRates(ctx context.Context) ([]entities.AveragePrivateUserRate, error)
}

//
// summaryRepository Implementation
//

type summaryRepository struct {
	logger         logger.Logger
	postgreQueries *rdb.PostgreQueries
}

func NewSummaryRepository(
	logger logger.Logger,
	postgreQueries *rdb.PostgreQueries,
) *summaryRepository {
	return &summaryRepository{
		logger:         logger,
		postgreQueries: postgreQueries,
	}
}

func (s *summaryRepository) Close(ctx context.Context) {
	s.postgreQueries.Close(ctx)
}

// PostgreSQL

func (s *summaryRepository) GetAllURLs(ctx context.Context) ([]entities.URL, error) {
	return s.postgreQueries.GetAllURLs(ctx)
}

func (s *summaryRepository) GetURLsByURLAddresses(
	ctx context.Context,
	urls []string,
) ([]entities.URL, error) {
	return s.postgreQueries.GetURLsByURLAddresses(ctx, urls)
}

// func (s *summaryRepository) GetUsersByURL(ctx context.Context, url string) ([]entities.RDBUser, error) {
// 	return s.postgreQueries.GetUsersByURL(ctx, url)
// }

func (s *summaryRepository) GetAveragePrivateUserRates(
	ctx context.Context,
) ([]entities.AveragePrivateUserRate, error) {
	return s.postgreQueries.GetAveragePrivateUserRates(ctx)
}
