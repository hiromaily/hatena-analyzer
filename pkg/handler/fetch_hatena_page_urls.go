package handler

import (
	"context"

	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/usecase"
)

//
// fetchURLsCLIHandler
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
