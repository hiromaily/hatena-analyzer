package usecase

import (
	"context"
	"fmt"

	"github.com/hiromaily/hatena-analyzer/pkg/entities"
	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/repository"
	"github.com/hiromaily/hatena-analyzer/pkg/tracer"
)

type ViewSummaryUsecaser interface {
	Execute(ctx context.Context, urls []string, threshold uint) error
}

type summaryUsecase struct {
	logger      logger.Logger
	tracer      tracer.Tracer
	summaryRepo repository.SummaryRepositorier
}

func NewViewSummaryUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	summaryRepo repository.SummaryRepositorier,
) (*summaryUsecase, error) {
	// validation
	// if len(urls) == 0 {
	// 	return nil, errors.New("urls is empty")
	// }

	return &summaryUsecase{
		logger:      logger,
		tracer:      tracer,
		summaryRepo: summaryRepo,
	}, nil
}

func (s *summaryUsecase) Execute(ctx context.Context, urls []string, threshold uint) error {
	// must be closed dbClient
	defer s.summaryRepo.Close(ctx)

	_, span := s.tracer.NewSpan(ctx, "summaryUsecase:Execute()")
	defer func() {
		span.End()
		s.tracer.Close(ctx)
	}()

	// get urls from DB if needed
	var entityURLs []entities.URL
	var err error
	if len(urls) == 0 {
		entityURLs, err = s.summaryRepo.GetAllURLs(ctx)
		if err != nil {
			s.logger.Error("failed to call bookmarkRepo.GetAllURLs()", "error", err)
			return err
		}
	} else {
		entityURLs, err = s.summaryRepo.GetURLsByURLAddresses(ctx, urls)
		if err != nil {
			s.logger.Error(
				"failed to call bookmarkDetailsRepo.GetURLsByURLAddresses()",
				"url_count", len(urls),
				"error", err,
			)
			return err
		}
	}

	s.logger.Info("url count", "count", len(entityURLs))

	fmt.Printf("[Private user rate over threshold: %d]\n", threshold)
	for _, entityURL := range entityURLs {
		if entityURL.PrivateUserRate > float64(threshold) {
			s.logger.Info(
				"url info",
				"url", entityURL.Address,
				"title", entityURL.Title,
				"bm_count", entityURL.BookmarkCount,
				"user_count", entityURL.NamedUserCount,
				"private_user_rate", entityURL.PrivateUserRate,
			)
		}
	}
	fmt.Println("")

	averageRates, err := s.summaryRepo.GetAveragePrivateUserRates(ctx)
	if err != nil {
		s.logger.Error("failed to call summaryRepo.GetAveragePrivateUserRates()", "error", err)
		return err
	}

	fmt.Println("[Average private user rate per category]")
	for _, ave := range averageRates {
		s.logger.Info(
			"average private user rate",
			"ave", ave.CategoryCode.String(),
			"average_private_user_rate", ave.AveragePrivateUserRate,
		)
	}

	return nil
}
