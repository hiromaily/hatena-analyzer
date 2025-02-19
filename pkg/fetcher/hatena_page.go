package fetcher

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
)

type hatenaPageFetcher struct {
	logger logger.Logger
}

func NewHatenaPageFetcher(logger logger.Logger) *hatenaPageFetcher {
	return &hatenaPageFetcher{
		logger: logger,
	}
}

// Fetch bookmark count of user from Hatena user's page

func (h *hatenaPageFetcher) Fetch(_ context.Context, url string) ([]string, error) {
	h.logger.Debug("fetching page: ", "url", url)
	return nil, nil
}
