package usecase

import (
	"context"
	"fmt"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
)

type FetchHatenaPageURLsUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchHatenaPageURLsUsecase struct {
	logger               logger.Logger
	tracer               tracer.Tracer
	urlRepo              repository.URLRepositorier
	hatenaPageURLFetcher fetcher.HatenaPageURLFetcher
	categoryCode         entities.CategoryCode
}

// TODO
// - add cli parameter: category_code

func NewFetchHatenaPageURLsUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	urlRepo repository.URLRepositorier,
	hatenaPageURLFetcher fetcher.HatenaPageURLFetcher,
	categoryCode entities.CategoryCode,
) (*fetchHatenaPageURLsUsecase, error) {
	// validation

	return &fetchHatenaPageURLsUsecase{
		logger:               logger,
		tracer:               tracer,
		urlRepo:              urlRepo,
		hatenaPageURLFetcher: hatenaPageURLFetcher,
		categoryCode:         categoryCode,
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

	// targetURLs := []string{
	// 	//"https://b.hatena.ne.jp/hotentry/all",
	// 	//"https://b.hatena.ne.jp/entrylist/all",
	// 	//"https://b.hatena.ne.jp/hotentry/general",
	// 	"https://b.hatena.ne.jp/hotentry/it",
	// }

	targetURLs := []string{}
	if f.categoryCode == entities.Unknown {
		categoryCodes := entities.GetCategoryCodeList()
		for _, code := range categoryCodes {
			targetURLs = append(
				targetURLs,
				fmt.Sprintf("%s/%s", "https://b.hatena.ne.jp/hotentry", code.String()),
			)
		}
	} else {
		targetURLs = append(targetURLs, fmt.Sprintf("%s/%s", "https://b.hatena.ne.jp/hotentry", f.categoryCode.String()))
	}

	totalFetchedURLs := []string{}
	for _, url := range targetURLs {
		// fetch page
		pageURLs, err := f.hatenaPageURLFetcher.Fetch(ctx, url)
		if err != nil {
			f.logger.Error("failed to fetch page", "url", url, "error", err)
			return err
		}
		if len(pageURLs) == 0 {
			f.logger.Warn("no URLs are fetched", "url", url)
			continue
		}
		category, err := entities.ExtractCategoryFromURL(url)
		if err != nil {
			f.logger.Error("failed to extract category from URL", "url", url, "error", err)
			continue
		}

		// Insert fetched URLs to DB
		// FIXME: duplicate key value violates unique constraint "urls_url_address_key" (SQLSTATE 23505)
		// TODO: create stored procedure to avoid conflict error
		if err := f.urlRepo.InsertURLs(ctx, category, totalFetchedURLs); err != nil {
			f.logger.Error("failed to insert URLs", "error", err)
		}
		totalFetchedURLs = append(totalFetchedURLs, pageURLs...)
	}
	f.logger.Info("total fetched URLs", "count", len(totalFetchedURLs))

	return nil
}
