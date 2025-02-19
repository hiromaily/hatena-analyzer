package handler

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

//
// fetchURLsCLIHandler
//

type fetchURLsCLIHandler struct {
	logger  logger.Logger
	usecase usecase.FetchURLsUsecaser
}

func NewFetchURLsCLIHandler(
	logger logger.Logger,
	usecase usecase.FetchURLsUsecaser,
) *fetchURLsCLIHandler {
	return &fetchURLsCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (f *fetchURLsCLIHandler) Handler(ctx context.Context) error {
	f.logger.Info("fetchURLsCLIHandler Handler")

	err := f.usecase.Execute(ctx)
	if err != nil {
		f.logger.Error("failed to fetch urls from page", "error", err)
	}
	return err
}
