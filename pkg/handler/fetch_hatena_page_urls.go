package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/usecase"
)

//
// fetchHatenaPageURLsCLIHandler
//

type fetchHatenaPageURLsCLIHandler struct {
	logger  logger.Logger
	usecase usecase.FetchHatenaPageURLsUsecaser
}

func NewFetchHatenaPageURLsCLIHandler(
	logger logger.Logger,
	usecase usecase.FetchHatenaPageURLsUsecaser,
) *fetchHatenaPageURLsCLIHandler {
	return &fetchHatenaPageURLsCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (f *fetchHatenaPageURLsCLIHandler) Handler(ctx context.Context) error {
	f.logger.Info("fetchHatenaPageURLsCLIHandler Handler")

	err := f.usecase.Execute(ctx)
	if err != nil {
		f.logger.Error("failed to fetch urls from page", "error", err)
	}
	return err
}

// dummy
func (f *fetchHatenaPageURLsCLIHandler) WebHandler(_ *gin.Context) {
}

//
// fetchHatenaPageURLsWebHandler
//

type fetchHatenaPageURLsWebHandler struct {
	logger  logger.Logger
	usecase usecase.FetchHatenaPageURLsUsecaser
}

func NewFetchHatenaPageURLsWebHandler(
	logger logger.Logger,
	usecase usecase.FetchHatenaPageURLsUsecaser,
) *fetchHatenaPageURLsWebHandler {
	return &fetchHatenaPageURLsWebHandler{
		logger:  logger,
		usecase: usecase,
	}
}

// dummy
func (f *fetchHatenaPageURLsWebHandler) Handler(_ context.Context) error {
	return nil
}

func (f *fetchHatenaPageURLsWebHandler) WebHandler(c *gin.Context) {
	ctx := c.Request.Context()

	err := f.usecase.Execute(ctx)
	if err != nil {
		f.logger.Error("failed to fetch bookmark data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookmark data"})
		return
	}

	f.logger.Info("successfully fetched bookmark data")
	c.JSON(http.StatusOK, gin.H{"message": "successfully fetched bookmark data"})
}
