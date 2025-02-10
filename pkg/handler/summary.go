package handler

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

//
// summaryCLIHandler
//

type summaryCLIHandler struct {
	logger  logger.Logger
	usecase usecase.ViewSummaryUsecaser
}

func NewViewSummaryCLIHandler(logger logger.Logger, usecase usecase.ViewSummaryUsecaser) *summaryCLIHandler {
	return &summaryCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (s *summaryCLIHandler) Handler(ctx context.Context) error {
	s.logger.Info("summaryCLIHandler Handler")

	err := s.usecase.Execute(ctx)
	if err != nil {
		s.logger.Error("failed to view bookmark summary data", "error", err)
	}
	return err
}
