package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/times"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
)

type ViewTimeSeriesUsecaser interface {
	Execute(ctx context.Context) error
}

type timeSeriesUsecase struct {
	logger         logger.Logger
	tracer         tracer.Tracer
	timeSeriesRepo repository.TimeSeriesRepositorier
	urls           []string
}

func NewViewTimeSeriesUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	timeSeriesRepo repository.TimeSeriesRepositorier,
	urls []string,
) (*timeSeriesUsecase, error) {
	// validation
	if len(urls) == 0 {
		return nil, errors.New("urls is empty")
	}

	return &timeSeriesUsecase{
		logger:         logger,
		tracer:         tracer,
		timeSeriesRepo: timeSeriesRepo,
		urls:           urls,
	}, nil
}

func (s *timeSeriesUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer s.timeSeriesRepo.Close(ctx)

	_, span := s.tracer.NewSpan(ctx, "timeSeriesUsecase:Execute()")
	defer func() {
		span.End()
		s.tracer.Close(ctx)
	}()

	for _, url := range s.urls {
		// get summaries from InfluxDB
		summaries, err := s.timeSeriesRepo.ReadEntitySummaries(ctx, url)
		if err != nil {
			s.logger.Error("failed to call summaryRepo.ReadEntitySummaries()", "url", url, "error", err)
			continue
		}
		if len(summaries) == 0 {
			s.logger.Warn("no data", "url", url)
			continue
		}

		fmt.Println("----------------------------------------------------------------------")
		fmt.Printf(" Title: %s,\n URL: %s\n", summaries[0].Title, url)
		fmt.Printf(" Time series\n")
		for _, summary := range summaries {
			fmt.Printf(
				"  - %s: total_bookmark: %d, user_count: %d, deleted_user_count: %d, private user rate: %.1f\n",
				times.ToJPTime(summary.Timestamp).Format(time.RFC3339),
				summary.Count,
				summary.UserCount,
				summary.DeletedUserCount,
				entities.PrivateUserRate(summary.Count, summary.UserCount),
			)
		}
	}

	return nil
}
