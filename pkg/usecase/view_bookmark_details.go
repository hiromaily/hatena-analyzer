package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
)

type ViewBookmarkDetailsUsecaser interface {
	Execute(ctx context.Context) error
}

type bookmarkDetailsUsecase struct {
	logger              logger.Logger
	tracer              tracer.Tracer
	bookmarkDetailsRepo repository.BookmarkDetailsRepositorier
	urls                []string
}

func NewViewBookmarkDetailsUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	bookmarkDetailsRepo repository.BookmarkDetailsRepositorier,
	urls []string,
) (*bookmarkDetailsUsecase, error) {
	// validation
	if len(urls) == 0 {
		return nil, errors.New("urls is empty")
	}

	return &bookmarkDetailsUsecase{
		logger:              logger,
		tracer:              tracer,
		bookmarkDetailsRepo: bookmarkDetailsRepo,
		urls:                urls,
	}, nil
}

func (b *bookmarkDetailsUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer b.bookmarkDetailsRepo.Close(ctx)

	_, span := b.tracer.NewSpan(ctx, "bookmarkDetailsUsecase:Execute()")
	defer func() {
		span.End()
		b.tracer.Close(ctx)
	}()

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

	for _, url := range b.urls {
		// get user by URL info from DB
		users, err := b.bookmarkDetailsRepo.GetUsersByURL(ctx, url)
		if err != nil {
			b.logger.Error("failed to call bookmarkDetailsRepo.GetUsersByURL()", "url", url, "error", err)
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
		newUserRate := float64(count10) / float64(9999) * 100

		fmt.Printf(" User's bookmark count / number of users whose bookmark count \n")
		fmt.Printf("  less 10:      %5d\n", count10)
		fmt.Printf("  less 100:     %5d\n", count100)
		fmt.Printf("  less 1000:    %5d\n", count1000)
		fmt.Printf("  less 10000:   %5d\n", count10000)
		fmt.Printf("  over 10000:   %5d\n", countOver)
		fmt.Printf(" New user rate:  %.1f\n", newUserRate)
	}

	return nil
}
