package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/repository"
	"github.com/hiromaily/hatena-analyzer/pkg/tracer"
)

type ViewBookmarkDetailsUsecaser interface {
	Execute(ctx context.Context, urls []string) error
}

type bookmarkDetailsUsecase struct {
	logger              logger.Logger
	tracer              tracer.Tracer
	bookmarkDetailsRepo repository.BookmarkDetailsRepositorier
}

func NewViewBookmarkDetailsUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	bookmarkDetailsRepo repository.BookmarkDetailsRepositorier,
) (*bookmarkDetailsUsecase, error) {
	// validation

	return &bookmarkDetailsUsecase{
		logger:              logger,
		tracer:              tracer,
		bookmarkDetailsRepo: bookmarkDetailsRepo,
	}, nil
}

func (b *bookmarkDetailsUsecase) Execute(ctx context.Context, urls []string) error {
	b.logger.Info("bookmarkDetailsUsecase Execute", "urls length", len(urls))

	// must be closed dbClient
	// defer b.bookmarkDetailsRepo.Close(ctx)

	_, span := b.tracer.NewSpan(ctx, "bookmarkDetailsUsecase:Execute()")
	defer func() {
		span.End()
		b.tracer.Close(ctx)
	}()

	// validation
	if len(urls) == 0 {
		return errors.New("urls is empty")
	}

	// get urls from DB if needed
	// if len(b.urls) == 0 {
	// 	var err error
	// 	entityURLs, err := b.bookmarkDetailsRepo.GetAllURLs(ctx)
	// 	if err != nil {
	// 		b.logger.Error("failed to call bookmarkDetailsRepo.GetAllURLs()", "error", err)
	// 		return err
	// 	}
	// 	b.urls = entities.FilterURLAddress(entityURLs)
	// }

	// get url info from DB
	urlModels, err := b.bookmarkDetailsRepo.GetURLsByURLAddresses(ctx, urls)
	if err != nil {
		b.logger.Error(
			"failed to call bookmarkDetailsRepo.GetURLsByURLAddresses()",
			"url_count", len(urls),
			"error", err,
		)
		return err
	}

	for _, urlModel := range urlModels {
		// get user by URL info from DB
		users, err := b.bookmarkDetailsRepo.GetUsersByURL(ctx, urlModel.Address)
		if err != nil {
			b.logger.Error(
				"failed to call bookmarkDetailsRepo.GetUsersByURL()",
				"url", urlModel.Address,
				"error", err,
			)
			continue
		}
		var count10, count100, count1000, count10000, countOver int
		for _, user := range users {
			switch {
			case user.BookmarkCount < 10:
				count10++
			case user.BookmarkCount < 100:
				count100++
			case user.BookmarkCount < 1000:
				count1000++
			case user.BookmarkCount < 10000:
				count10000++
			default:
				countOver++
			}
		}
		// calculate average
		// less 10 user must be suspicious
		newUserRate := float64(count10) / float64(urlModel.NamedUserCount) * 100

		fmt.Println("----------------------------------------------------------------------")
		fmt.Printf(" Title: %s,\n URL: %s\n", urlModel.Title, urlModel.Address)
		fmt.Printf(" User's bookmark count / number of users whose bookmark count \n")
		fmt.Printf(" - less 10:      %5d\n", count10)
		fmt.Printf(" - less 100:     %5d\n", count100)
		fmt.Printf(" - less 1000:    %5d\n", count1000)
		fmt.Printf(" - less 10000:   %5d\n", count10000)
		fmt.Printf(" - over 10000:   %5d\n", countOver)
		fmt.Printf(" New user rate:  %.1f\n", newUserRate)
	}

	return nil
}
