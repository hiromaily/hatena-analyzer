package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/usecase"
)

//
// viewSummaryCLIHandler
//

type viewSummaryCLIHandler struct {
	logger    logger.Logger
	usecase   usecase.ViewSummaryUsecaser
	urls      []string
	threshold uint
}

func NewViewSummaryCLIHandler(
	logger logger.Logger,
	usecase usecase.ViewSummaryUsecaser,
	urls []string,
	threshold uint,
) *viewSummaryCLIHandler {
	if threshold == 0 {
		// default
		threshold = 50
	}

	return &viewSummaryCLIHandler{
		logger:    logger,
		usecase:   usecase,
		urls:      urls,
		threshold: threshold,
	}
}

func (v *viewSummaryCLIHandler) Handler(ctx context.Context) error {
	v.logger.Info("viewSummaryCLIHandler Handler")

	err := v.usecase.Execute(ctx, v.urls, v.threshold)
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

	// request
	urlString := c.DefaultQuery("urls", "")
	var urls []string
	if urlString != "" {
		urls = strings.Split(urlString, ",")
		v.logger.Info("given URLs", "urls", urls, "len", len(urls))
	}

	// threshold := c.Query("threshold")
	err := v.usecase.Execute(ctx, urls, 50)
	if err != nil {
		v.logger.Error("failed to fetch bookmark data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookmark data"})
		return
	}

	v.logger.Info("successfully fetched bookmark data")
	c.JSON(http.StatusOK, gin.H{"message": "successfully fetched bookmark data"})
}
