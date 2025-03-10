package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hiromaily/hatena-analyzer/pkg/entities"
	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/repository"
	"github.com/hiromaily/hatena-analyzer/pkg/times"
	"github.com/hiromaily/hatena-analyzer/pkg/tracer"
)

type ViewTimeSeriesUsecaser interface {
	Execute(ctx context.Context, urls []string) error
}

type timeSeriesUsecase struct {
	logger         logger.Logger
	tracer         tracer.Tracer
	timeSeriesRepo repository.TimeSeriesRepositorier
}

func NewViewTimeSeriesUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	timeSeriesRepo repository.TimeSeriesRepositorier,
) (*timeSeriesUsecase, error) {
	return &timeSeriesUsecase{
		logger:         logger,
		tracer:         tracer,
		timeSeriesRepo: timeSeriesRepo,
	}, nil
}

func (t *timeSeriesUsecase) Execute(ctx context.Context, urls []string) error {
	t.logger.Info("timeSeriesUsecase Execute", "urls length", len(urls))

	// must be closed dbClient
	// defer t.timeSeriesRepo.Close(ctx)

	_, span := t.tracer.NewSpan(ctx, "timeSeriesUsecase:Execute()")
	defer func() {
		span.End()
		t.tracer.Close(ctx)
	}()

	// validation
	if len(urls) == 0 {
		return errors.New("urls is empty")
	}

	for _, url := range urls {
		// get summaries from InfluxDB
		summaries, err := t.timeSeriesRepo.ReadEntitySummaries(ctx, url)
		if err != nil {
			t.logger.Error("failed to call timeSeriesRepo.ReadEntitySummaries()", "url", url, "error", err)
			continue
		}
		if len(summaries) == 0 {
			t.logger.Warn("no data", "url", url)
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
