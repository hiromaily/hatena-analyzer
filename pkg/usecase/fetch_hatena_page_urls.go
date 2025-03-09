package usecase

import (
	"context"
	"fmt"

	"github.com/hiromaily/hatena-analyzer/pkg/entities"
	"github.com/hiromaily/hatena-analyzer/pkg/fetcher"
	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/repository"
	"github.com/hiromaily/hatena-analyzer/pkg/tracer"
)

type FetchHatenaPageURLsUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchHatenaPageURLsUsecase struct {
	logger               logger.Logger
	tracer               tracer.Tracer
	fetchURLRepo         repository.FetchURLRepositorier
	hatenaPageURLFetcher fetcher.HatenaPageURLFetcher
	categoryCode         entities.CategoryCode
}

// TODO
// - fetch urls concurrently
// - add test

func NewFetchHatenaPageURLsUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	fetchURLRepo repository.FetchURLRepositorier,
	hatenaPageURLFetcher fetcher.HatenaPageURLFetcher,
	categoryCode entities.CategoryCode,
) (*fetchHatenaPageURLsUsecase, error) {
	// validation

	return &fetchHatenaPageURLsUsecase{
		logger:               logger,
		tracer:               tracer,
		fetchURLRepo:         fetchURLRepo,
		hatenaPageURLFetcher: hatenaPageURLFetcher,
		categoryCode:         categoryCode,
	}, nil
}

// Fetch bookmark users, title, count related given URLs using Hatena entity API and save data to DB

func (f *fetchHatenaPageURLsUsecase) Execute(ctx context.Context) error {
	f.logger.Info("fetchHatenaPageURLsUsecase Execute")

	// must be closed dbClient
	// defer f.fetchURLRepo.Close(ctx)

	_, span := f.tracer.NewSpan(ctx, "fetchURLsUsecase:Execute()")
	defer func() {
		span.End()
		f.tracer.Close(ctx)
	}()

	targetURLs := []string{}
	if f.categoryCode == entities.Unknown {
		// fetch all categories except `all: 総合`
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
		category, err := entities.ExtractCategoryFromURL(url)
		if err != nil {
			f.logger.Error("failed to extract category from URL", "url", url, "error", err)
			continue
		}

		// fetch page
		f.logger.Info("fetching page", "url", url)
		linkInfos, err := f.hatenaPageURLFetcher.Fetch(ctx, url, category == entities.All)
		if err != nil {
			f.logger.Error("failed to fetch page", "url", url, "error", err)
			return err
		}
		if len(linkInfos) == 0 {
			f.logger.Warn("no URLs are fetched", "url", url)
			continue
		}

		// Insert fetched URLs to DB
		// FIXED: duplicate key value violates unique constraint "urls_url_address_key" (SQLSTATE 23505)
		f.logger.Info("insert urls", "category", category.String(), "url_count", len(linkInfos))
		// if err := f.fetchURLRepo.InsertURLs(ctx, category, pageURLs); err != nil {
		// 	f.logger.Error("failed to insert URLs", "category", category.String(), "error", err)
		// }
		urls, categories, isAlls := entities.LinkInfos(linkInfos).Extract()
		if err := f.fetchURLRepo.CallBulkInsertURLs(ctx, urls, categories, isAlls); err != nil {
			f.logger.Error("failed to insert URLs", "category", category.String(), "error", err)
		}
		totalFetchedURLs = append(totalFetchedURLs, urls...)
	}
	f.logger.Info("total fetched URLs", "total_url_count", len(totalFetchedURLs))

	return nil
}
