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
// viewTimeSeriesCLIHandler
//

type viewTimeSeriesCLIHandler struct {
	logger  logger.Logger
	usecase usecase.ViewTimeSeriesUsecaser
	urls    []string
}

func NewViewTimeSeriesCLIHandler(
	logger logger.Logger,
	usecase usecase.ViewTimeSeriesUsecaser,
	urls []string,
) *viewTimeSeriesCLIHandler {
	return &viewTimeSeriesCLIHandler{
		logger:  logger,
		usecase: usecase,
		urls:    urls,
	}
}

func (v *viewTimeSeriesCLIHandler) Handler(ctx context.Context) error {
	v.logger.Info("viewTimeSeriesCLIHandler Handler")

	err := v.usecase.Execute(ctx, v.urls)
	if err != nil {
		v.logger.Error("failed to view bookmark time series", "error", err)
	}
	return err
}

// dummy
func (v *viewTimeSeriesCLIHandler) WebHandler(_ *gin.Context) {
}

//
// viewTimeSeriesWebHandler
//

type viewTimeSeriesWebHandler struct {
	logger  logger.Logger
	usecase usecase.ViewTimeSeriesUsecaser
}

func NewViewTimeSeriesWebHandler(
	logger logger.Logger,
	usecase usecase.ViewTimeSeriesUsecaser,
) *viewTimeSeriesWebHandler {
	return &viewTimeSeriesWebHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (v *viewTimeSeriesWebHandler) Handler(_ context.Context) error {
	return nil
}

func (v *viewTimeSeriesWebHandler) WebHandler(c *gin.Context) {
	v.logger.Info("viewTimeSeriesWebHandler WebHandler")

	ctx := c.Request.Context()

	// request
	urlString := c.DefaultQuery("urls", "")
	var urls []string
	if urlString != "" {
		urls = strings.Split(urlString, ",")
		v.logger.Info("given URLs", "urls", urls, "len", len(urls))
	}

	err := v.usecase.Execute(ctx, urls)
	if err != nil {
		v.logger.Error("failed to fetch bookmark data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookmark data"})
		return
	}

	v.logger.Info("successfully fetched bookmark data")
	c.JSON(http.StatusOK, gin.H{"message": "successfully fetched bookmark data"})
}
