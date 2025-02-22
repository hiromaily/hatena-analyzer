package handler

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

//
// viewTimeSeriesCLIHandler
//

type viewTimeSeriesCLIHandler struct {
	logger  logger.Logger
	usecase usecase.ViewTimeSeriesUsecaser
}

func NewViewTimeSeriesCLIHandler(
	logger logger.Logger,
	usecase usecase.ViewSummaryUsecaser,
) *viewTimeSeriesCLIHandler {
	return &viewTimeSeriesCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (s *viewTimeSeriesCLIHandler) Handler(ctx context.Context) error {
	s.logger.Info("viewTimeSeriesCLIHandler Handler")

	err := s.usecase.Execute(ctx)
	if err != nil {
		s.logger.Error("failed to view bookmark time series", "error", err)
	}
	return err
}
