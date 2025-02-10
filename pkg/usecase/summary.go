package usecase

import (
	"context"
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
) *summaryUsecase {

	// target URL list
	urls := []string{
		"https://note.com/simplearchitect/n/nadc0bcdd5b3d",
		"https://note.com/simplearchitect/n/n871f29ffbfac",
		"https://note.com/simplearchitect/n/n86a95bc19b4c",
		"https://note.com/simplearchitect/n/nfd147540e3db",
	}

	return &summaryUsecase{
		logger:      logger,
		summaryRepo: summaryRepo,
		urls:        urls,
	}
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
				" %s: count: %d, user_count: %d, deleted_user_count: %d\n",
				times.ToJPTime(summary.Timestamp).Format(time.RFC3339),
				summary.Count,
				summary.UserCount,
				summary.DeletedUserCount,
			)
		}
	}

	return nil
}
