package handler

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

//
// viewSummaryCLIHandler
//

type viewSummaryCLIHandler struct {
	logger  logger.Logger
	usecase usecase.ViewSummaryUsecaser
}

func NewViewSummaryCLIHandler(
	logger logger.Logger,
	usecase usecase.ViewSummaryUsecaser,
) *viewSummaryCLIHandler {
	return &viewSummaryCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (s *viewSummaryCLIHandler) Handler(ctx context.Context) error {
	s.logger.Info("viewSummaryCLIHandler Handler")

	err := s.usecase.Execute(ctx)
	if err != nil {
		s.logger.Error("failed to view bookmark summary data", "error", err)
	}
	return err
}
