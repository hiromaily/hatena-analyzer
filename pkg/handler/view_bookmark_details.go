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
// viewBookmarkDetailsCLIHandler
//

type viewBookmarkDetailsCLIHandler struct {
	logger  logger.Logger
	usecase usecase.ViewBookmarkDetailsUsecaser
	urls    []string
}

func NewViewBookmarkDetailsCLIHandler(
	logger logger.Logger,
	usecase usecase.ViewBookmarkDetailsUsecaser,
	urls []string,
) *viewBookmarkDetailsCLIHandler {
	return &viewBookmarkDetailsCLIHandler{
		logger:  logger,
		usecase: usecase,
		urls:    urls,
	}
}

func (v *viewBookmarkDetailsCLIHandler) Handler(ctx context.Context) error {
	v.logger.Info("viewBookmarkDetailsCLIHandler Handler")

	err := v.usecase.Execute(ctx, v.urls)
	if err != nil {
		v.logger.Error("failed to view bookmark details", "error", err)
	}
	return err
}

// dummy
func (v *viewBookmarkDetailsCLIHandler) WebHandler(_ *gin.Context) {
}

//
// viewBookmarkDetailsWebHandler
//

type viewBookmarkDetailsWebHandler struct {
	logger  logger.Logger
	usecase usecase.ViewBookmarkDetailsUsecaser
}

func NewViewBookmarkDetailsWebHandler(
	logger logger.Logger,
	usecase usecase.ViewBookmarkDetailsUsecaser,
) *viewBookmarkDetailsWebHandler {
	return &viewBookmarkDetailsWebHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (v *viewBookmarkDetailsWebHandler) Handler(_ context.Context) error {
	return nil
}

func (v *viewBookmarkDetailsWebHandler) WebHandler(c *gin.Context) {
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
