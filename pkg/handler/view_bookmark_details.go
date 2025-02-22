package handler

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

//
// viewBookmarkDetailsCLIHandler
//

type viewBookmarkDetailsCLIHandler struct {
	logger  logger.Logger
	usecase usecase.ViewBookmarkDetailsUsecaser
}

func NewViewBookmarkDetailsCLIHandler(
	logger logger.Logger,
	usecase usecase.ViewBookmarkDetailsUsecaser,
) *viewBookmarkDetailsCLIHandler {
	return &viewBookmarkDetailsCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (v *viewBookmarkDetailsCLIHandler) Handler(ctx context.Context) error {
	v.logger.Info("viewBookmarkDetailsCLIHandler Handler")

	err := v.usecase.Execute(ctx)
	if err != nil {
		v.logger.Error("failed to view bookmark details", "error", err)
	}
	return err
}
