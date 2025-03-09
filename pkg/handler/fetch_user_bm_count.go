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
// fetchUserBookmarkCountCLIHandler
//

type fetchUserBookmarkCountCLIHandler struct {
	logger  logger.Logger
	usecase usecase.FetchUserBookmarkCountUsecaser
	urls    []string
}

func NewFetchUserBookmarkCountCLIHandler(
	logger logger.Logger,
	usecase usecase.FetchUserBookmarkCountUsecaser,
	urls []string,
) *fetchUserBookmarkCountCLIHandler {
	return &fetchUserBookmarkCountCLIHandler{
		logger:  logger,
		usecase: usecase,
		urls:    urls,
	}
}

func (f *fetchUserBookmarkCountCLIHandler) Handler(ctx context.Context) error {
	f.logger.Info("fetchUserBookmarkCountCLIHandler Handler")

	err := f.usecase.Execute(ctx, f.urls)
	if err != nil {
		f.logger.Error("failed to update user info", "error", err)
	}
	return err
}

// dummy
func (f *fetchUserBookmarkCountCLIHandler) WebHandler(_ *gin.Context) {
}

//
// fetchUserBookmarkCountWebHandler
//

type fetchUserBookmarkCountWebHandler struct {
	logger  logger.Logger
	usecase usecase.FetchUserBookmarkCountUsecaser
}

func NewFetchUserBookmarkCountWebHandler(
	logger logger.Logger,
	usecase usecase.FetchUserBookmarkCountUsecaser,
) *fetchUserBookmarkCountWebHandler {
	return &fetchUserBookmarkCountWebHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (f *fetchUserBookmarkCountWebHandler) Handler(_ context.Context) error {
	return nil
}

// type fetchUserBookmarkCountRequest struct {
// 	URLs []string `json:"urls"`
// }

func (f *fetchUserBookmarkCountWebHandler) WebHandler(c *gin.Context) {
	f.logger.Info("fetchUserBookmarkCountWebHandler WebHandler")

	ctx := c.Request.Context()

	// request
	urlString := c.DefaultQuery("urls", "")
	var urls []string
	if urlString != "" {
		urls = strings.Split(urlString, ",")
		f.logger.Info("given URLs", "urls", urls, "len", len(urls))
	}

	err := f.usecase.Execute(ctx, urls)
	if err != nil {
		f.logger.Error("failed to fetch bookmark data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookmark data"})
		return
	}

	f.logger.Info("successfully fetched bookmark data")
	c.JSON(http.StatusOK, gin.H{"message": "successfully fetched bookmark data"})
}
