package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

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

// dummy
func (v *viewSummaryCLIHandler) WebHandler(_ *gin.Context) {
}

//
// viewSummaryWebHandler
//

type viewSummaryWebHandler struct {
	logger  logger.Logger
	usecase usecase.ViewSummaryUsecaser
}

func NewViewSummaryWebHandler(
	logger logger.Logger,
	usecase usecase.ViewSummaryUsecaser,
) *viewSummaryWebHandler {
	return &viewSummaryWebHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (v *viewSummaryWebHandler) Handler(_ context.Context) error {
	return nil
}

func (v *viewSummaryWebHandler) WebHandler(c *gin.Context) {
	ctx := c.Request.Context()

	err := v.usecase.Execute(ctx)
	if err != nil {
		v.logger.Error("failed to fetch bookmark data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookmark data"})
		return
	}

	v.logger.Info("successfully fetched bookmark data")
	c.JSON(http.StatusOK, gin.H{"message": "successfully fetched bookmark data"})
}
