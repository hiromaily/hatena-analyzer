package usecase

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
)

type FetchHatenaPageURLsUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchHatenaPageURLsUsecase struct {
	logger            logger.Logger
	tracer            tracer.Tracer
	urlRepo           repository.URLRepositorier
	hatenaPageFetcher fetcher.HatenaPageFetcher
	targetURLs        []string
}

func NewFetchHatenaPageURLsUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	urlRepo repository.URLRepositorier,
	hatenaPageFetcher fetcher.HatenaPageFetcher,
) (*fetchHatenaPageURLsUsecase, error) {
	// validation

	targetURLs := []string{
		"https://b.hatena.ne.jp/hotentry/all",
		"https://b.hatena.ne.jp/entrylist/all",
	}

	return &fetchHatenaPageURLsUsecase{
		logger:            logger,
		tracer:            tracer,
		urlRepo:           urlRepo,
		hatenaPageFetcher: hatenaPageFetcher,
		targetURLs:        targetURLs,
	}, nil
}

// Fetch bookmark users, title, count related given URLs using Hatena entity API and save data to DB

func (f *fetchHatenaPageURLsUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer f.urlRepo.Close(ctx)

	_, span := f.tracer.NewSpan(ctx, "fetchURLsUsecase:Execute()")
	defer func() {
		span.End()
		f.tracer.Close(ctx)
	}()

	totalFetchedURLs := []string{}
	for _, url := range f.targetURLs {
		// fetch page
		pageURLs, err := f.hatenaPageFetcher.Fetch(ctx, url)
		if err != nil {
			f.logger.Error("failed to fetch page", "url", url, "error", err)
			return err
		}
		if len(pageURLs) == 0 {
			f.logger.Warn("no URLs are fetched", "url", url)
			continue
		}
		totalFetchedURLs = append(totalFetchedURLs, pageURLs...)
	}
	// save fetched URLs to DB
	if len(totalFetchedURLs) == 0 {
		f.logger.Warn("no URLs are fetched")
		return nil
	}
	return f.urlRepo.InsertURLs(ctx, totalFetchedURLs)
}
