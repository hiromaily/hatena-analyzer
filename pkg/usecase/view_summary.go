package usecase

import (
	"context"
	"fmt"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
)

type ViewSummaryUsecaser interface {
	Execute(ctx context.Context) error
}

type summaryUsecase struct {
	logger      logger.Logger
	tracer      tracer.Tracer
	summaryRepo repository.SummaryRepositorier
	urls        []string
	threshold   uint
}

func NewViewSummaryUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	summaryRepo repository.SummaryRepositorier,
	urls []string,
	threshold uint,
) (*summaryUsecase, error) {
	// validation
	// if len(urls) == 0 {
	// 	return nil, errors.New("urls is empty")
	// }

	return &summaryUsecase{
		logger:      logger,
		tracer:      tracer,
		summaryRepo: summaryRepo,
		urls:        urls,
		threshold:   threshold,
	}, nil
}

func (s *summaryUsecase) Execute(ctx context.Context) error {
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
	if len(s.urls) == 0 {
		entityURLs, err = s.summaryRepo.GetAllURLs(ctx)
		if err != nil {
			s.logger.Error("failed to call bookmarkRepo.GetAllURLs()", "error", err)
			return err
		}
	} else {
		entityURLs, err = s.summaryRepo.GetURLsByURLAddresses(ctx, s.urls)
		if err != nil {
			s.logger.Error(
				"failed to call bookmarkDetailsRepo.GetURLsByURLAddresses()",
				"url_count", len(s.urls),
				"error", err,
			)
			return err
		}
	}

	s.logger.Info("url count", "count", len(entityURLs))

	fmt.Printf("[Private user rate over threshold: %d]\n", s.threshold)
	for _, entityURL := range entityURLs {
		if entityURL.PrivateUserRate > float64(s.threshold) {
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
