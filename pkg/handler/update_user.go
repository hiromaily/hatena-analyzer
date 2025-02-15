package handler

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

//
// updateUserInfoCLIHandler
//

type updateUserInfoCLIHandler struct {
	logger  logger.Logger
	usecase usecase.UpdateUserInfoUsecaser
}

func NewUpdateUserInfoCLIHandler(
	logger logger.Logger,
	usecase usecase.UpdateUserInfoUsecaser,
) *updateUserInfoCLIHandler {
	return &updateUserInfoCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (f *updateUserInfoCLIHandler) Handler(ctx context.Context) error {
	f.logger.Info("updateUserInfoCLIHandler Handler")

	err := f.usecase.Execute(ctx)
	if err != nil {
		f.logger.Error("failed to update user info", "error", err)
	}
	return err
}
