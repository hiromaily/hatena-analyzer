package handler

import (
	"context"

	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/usecase"
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

func (v *viewSummaryCLIHandler) Handler(ctx context.Context) error {
	v.logger.Info("viewSummaryCLIHandler Handler")

	err := v.usecase.Execute(ctx)
	if err != nil {
		v.logger.Error("failed to view bookmark summary data", "error", err)
	}
	return err
}
