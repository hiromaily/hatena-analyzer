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
// fetchBookmarkCLIHandler
//

type fetchBookmarkCLIHandler struct {
	logger    logger.Logger
	usecase   usecase.FetchBookmarkUsecaser
	urls      []string
	isVerbose bool
}

func NewFetchBookmarkCLIHandler(
	logger logger.Logger,
	usecase usecase.FetchBookmarkUsecaser,
	urls []string,
	isVerbose bool,
) *fetchBookmarkCLIHandler {
	return &fetchBookmarkCLIHandler{
		logger:    logger,
		usecase:   usecase,
		urls:      urls,
		isVerbose: isVerbose,
	}
}

func (f *fetchBookmarkCLIHandler) Handler(ctx context.Context) error {
	f.logger.Info("fetchBookmarkCLIHandler Handler")

	err := f.usecase.Execute(ctx, f.urls, f.isVerbose)
	if err != nil {
		f.logger.Error("failed to fetch bookmark data", "error", err)
	}
	return err
}

// dummy
func (f *fetchBookmarkCLIHandler) WebHandler(_ *gin.Context) {
}

//
// fetchBookmarkWebHandler
//

type fetchBookmarkWebHandler struct {
	logger  logger.Logger
	usecase usecase.FetchBookmarkUsecaser
}

func NewFetchBookmarkWebHandler(
	logger logger.Logger,
	usecase usecase.FetchBookmarkUsecaser,
) *fetchBookmarkWebHandler {
	return &fetchBookmarkWebHandler{
		logger:  logger,
		usecase: usecase,
	}
}

// dummy
func (f *fetchBookmarkWebHandler) Handler(_ context.Context) error {
	return nil
}

func (f *fetchBookmarkWebHandler) WebHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// request
	urlString := c.DefaultQuery("urls", "")
	var urls []string
	if urlString != "" {
		urls = strings.Split(urlString, ",")
		f.logger.Info("given URLs", "urls", urls, "len", len(urls))
	}

	err := f.usecase.Execute(ctx, urls, false)
	if err != nil {
		f.logger.Error("failed to fetch bookmark data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookmark data"})
		return
	}

	f.logger.Info("successfully fetched bookmark data")
	c.JSON(http.StatusOK, gin.H{"message": "successfully fetched bookmark data"})
}
