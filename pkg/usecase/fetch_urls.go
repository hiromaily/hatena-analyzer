package usecase

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
)

type FetchURLsUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchURLsUsecase struct {
	logger            logger.Logger
	tracer            tracer.Tracer
	bookmarkRepo      repository.BookmarkRepositorier
	hatenaPageFetcher fetcher.HatenaPageFetcher
	targetURLs        []string
}

func NewFetchURLsUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	bookmarkRepo repository.BookmarkRepositorier,
	hatenaPageFetcher fetcher.HatenaPageFetcher,
) (*fetchURLsUsecase, error) {
	// validation

	targetURLs := []string{
		"https://b.hatena.ne.jp/hotentry/all",
		"https://b.hatena.ne.jp/entrylist/all",
	}

	return &fetchURLsUsecase{
		logger:            logger,
		tracer:            tracer,
		bookmarkRepo:      bookmarkRepo,
		hatenaPageFetcher: hatenaPageFetcher,
		targetURLs:        targetURLs,
	}, nil
}

// Fetch bookmark users, title, count related given URLs using Hatena entity API and save data to DB

func (f *fetchURLsUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer f.bookmarkRepo.Close(ctx)

	_, span := f.tracer.NewSpan(ctx, "fetchURLsUsecase:Execute()")
	defer func() {
		span.End()
		f.tracer.Close(ctx)
	}()

	// for _, url := range f.targetURLs {
	// 	// fetch page
	// }

	return nil
}
