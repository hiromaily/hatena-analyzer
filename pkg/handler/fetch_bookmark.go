package handler

import (
	"context"

	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/usecase"
)

//
// fetchBookmarkCLIHandler
//

type fetchBookmarkCLIHandler struct {
	logger  logger.Logger
	usecase usecase.FetchBookmarkUsecaser
}

func NewFetchBookmarkCLIHandler(
	logger logger.Logger,
	usecase usecase.FetchBookmarkUsecaser,
) *fetchBookmarkCLIHandler {
	return &fetchBookmarkCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (f *fetchBookmarkCLIHandler) Handler(ctx context.Context) error {
	f.logger.Info("fetchBookmarkCLIHandler Handler")

	err := f.usecase.Execute(ctx)
	if err != nil {
		f.logger.Error("failed to fetch bookmark data", "error", err)
	}
	return err
}
