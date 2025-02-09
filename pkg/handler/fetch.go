package handler

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

//
// fetchCLIHandler
//

type fetchCLIHandler struct {
	logger  logger.Logger
	usecase usecase.FetchUsecaser
}

func NewFetchCLIHandler(logger logger.Logger, usecase usecase.FetchUsecaser) *fetchCLIHandler {
	return &fetchCLIHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (f *fetchCLIHandler) Handler(ctx context.Context) error {
	return f.usecase.Execute(ctx)
}
