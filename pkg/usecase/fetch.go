package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
)

type FetchUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchUsecase struct {
	logger          logger.Logger
	bookmarkRepo    repository.BookmarkRepositorier
	bookmarkFetcher fetcher.BookmarkFetcher
	urls            []string
}

func NewFetchUsecase(
	logger logger.Logger,
	bookmarkRepo repository.BookmarkRepositorier,
	bookmarkFetcher fetcher.BookmarkFetcher,
	urls []string,
) (*fetchUsecase, error) {
	// validation
	if len(urls) == 0 {
		return nil, errors.New("urls is empty")
	}

	return &fetchUsecase{
		logger:          logger,
		bookmarkRepo:    bookmarkRepo,
		bookmarkFetcher: bookmarkFetcher,
		urls:            urls,
	}, nil
}

func (f *fetchUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer f.bookmarkRepo.Close(ctx)

	for _, url := range f.urls {
		// load bookmark summary
		bookmarkSummary, err := f.bookmarkRepo.ReadEntitySummary(ctx, url)
		if err != nil {
			f.logger.Error("failed to call bookmarkRepo.ReadEntitySummary()", "url", url, "error", err)
			continue
		}
		f.logger.Debug("bookmark summary loaded",
			"url", url,
			"bookmarkSummary.Title", bookmarkSummary.Title,
			"bookmarkSummary.Count", bookmarkSummary.Count,
			"bookmarkSummary.UserCount", bookmarkSummary.UserCount,
		)

		// load bookmark
		existingBookmark, err := f.bookmarkRepo.ReadEntity(ctx, url)
		if err != nil {
			f.logger.Error("failed to call bookmarkRepo.ReadEntity()", "url", url, "error", err)
			continue
		}

		if existingBookmark == nil {
			f.logger.Debug("data not found on MongoDB", "url", url)
			// initialize entities.Bookmark
			existingBookmark = &entities.Bookmark{}
			existingBookmark.Users = make(map[string]entities.User)
		}

		f.logger.Info("data loaded",
			"url", url,
			"existingBookmark.Title", existingBookmark.Title,
			"existingBookmark.Count", existingBookmark.Count,
			"existingBookmark.User.Length", len(existingBookmark.Users),
		)

		// 既存ユーザーをすべて`isDeleted = true`に設定
		for userName := range existingBookmark.Users {
			existingBookmark.Users[userName] = entities.User{
				Name:        userName,
				IsDeleted:   true,
				IsCommented: existingBookmark.Users[userName].IsCommented,
			}
		}

		// 新しいデータを取得
		newBookmark, err := f.bookmarkFetcher.Entity(ctx, url)
		if err != nil {
			f.logger.Error("failed to call fetchBookmarkData()", "url", url, "error", err)
			continue
		}
		f.logger.Info(
			"data fetched",
			"url", url,
			"newBookmark.Title", newBookmark.Title,
			"newBookmark.Count", newBookmark.Count,
			"newBookmark.User.Length", len(newBookmark.Users),
		)

		// 取得したユーザーで`isDeleted = false`に設定
		for userName, user := range newBookmark.Users {
			existingBookmark.Users[userName] = entities.User{
				Name:        userName,
				IsDeleted:   false,
				IsCommented: user.IsCommented,
			}
		}

		existingBookmark.Title = newBookmark.Title
		existingBookmark.Count = newBookmark.Count
		existingBookmark.Timestamp = newBookmark.Timestamp

		f.logger.Info("data will be stored",
			"url", url,
			"newBookmark.Title", existingBookmark.Title,
			"newBookmark.Count", existingBookmark.Count,
			"newBookmark.User.Length", len(existingBookmark.Users),
		)

		// save data
		err = f.bookmarkRepo.WriteEntitySummary(ctx, url, existingBookmark)
		if err != nil {
			f.logger.Error("failed to call bookmarkRepo.WriteEntitySummary()", "url", url, "error", err)
			continue
		}
		err = f.bookmarkRepo.WriteEntity(ctx, url, existingBookmark)
		if err != nil {
			f.logger.Error("failed to call bookmarkRepo.WriteEntity()", "url", url, "error", err)
			continue
		}
		f.logger.Info("data saved", "url", url)

		// Print data
		fmt.Println("===================================================================")
		fmt.Printf("Title: %s\n", existingBookmark.Title)
		fmt.Printf("Count: %d\n", existingBookmark.Count)
		fmt.Printf("UserCount: %d\n", len(existingBookmark.Users))
		fmt.Printf("DeletedUserCount: %d\n", existingBookmark.CountDeletedUser())

		// fmt.Printf("Users:\n")
		// for _, user := range existingBookmark.Users {
		// 	fmt.Printf(" - %s\n", user.Name)
		// }
		fmt.Println()
	}

	return nil
}
