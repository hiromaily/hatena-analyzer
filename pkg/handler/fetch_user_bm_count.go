package handler

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

//
// updateUserInfoCLIHandler
//

type fetchUserBookmarkCountCLIHandler struct {
	logger  logger.Logger
	usecase usecase.FetchUserBookmarkCountUsecaser
}

func NewFetchUserBookmarkCountCLIHandler(
	logger logger.Logger,
	usecase usecase.FetchUserBookmarkCountUsecaser,
) *fetchUserBookmarkCountCLIHandler {
	return &fetchUserBookmarkCountCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (f *fetchUserBookmarkCountCLIHandler) Handler(ctx context.Context) error {
	f.logger.Info("fetchUserBookmarkCountCLIHandler Handler")

	err := f.usecase.Execute(ctx)
	if err != nil {
		f.logger.Error("failed to update user info", "error", err)
	}
	return err
}
