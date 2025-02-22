package usecase

import (
	"context"

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
}

func NewViewSummaryUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	summaryRepo repository.SummaryRepositorier,
	urls []string,
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
	if len(s.urls) == 0 {
		var err error
		entityURLs, err := s.summaryRepo.GetAllURLs(ctx)
		if err != nil {
			s.logger.Error("failed to call bookmarkRepo.GetAllURLs()", "error", err)
			return err
		}
		s.urls = entities.FilterURLAddress(entityURLs)
	}

	return nil
}
