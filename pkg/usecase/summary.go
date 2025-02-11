package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/times"
)

type ViewSummaryUsecaser interface {
	Execute(ctx context.Context) error
}

type summaryUsecase struct {
	logger      logger.Logger
	summaryRepo repository.SummaryRepositorier
	urls        []string
}

func NewViewSummaryUsecase(
	logger logger.Logger,
	summaryRepo repository.SummaryRepositorier,
	urls []string,
) (*summaryUsecase, error) {
	// validation
	if len(urls) == 0 {
		return nil, errors.New("urls is empty")
	}

	return &summaryUsecase{
		logger:      logger,
		summaryRepo: summaryRepo,
		urls:        urls,
	}, nil
}

func (s *summaryUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer s.summaryRepo.Close(ctx)

	for _, url := range s.urls {
		// get summaries
		summaries, err := s.summaryRepo.ReadEntitySummaries(ctx, url)
		if err != nil {
			s.logger.Error("failed to call summaryRepo.ReadEntitySummaries()", "url", url, "error", err)
			continue
		}
		fmt.Printf("[Summary] URL: %s\n", url)

		for _, summary := range summaries {
			fmt.Printf(
				" %s: count: %d, user_count: %d, deleted_user_count: %d, delete rate: %.1f\n",
				times.ToJPTime(summary.Timestamp).Format(time.RFC3339),
				summary.Count,
				summary.UserCount,
				summary.DeletedUserCount,
				float64(summary.Count-summary.UserCount)/float64(summary.Count)*100,
			)
		}
	}

	return nil
}
