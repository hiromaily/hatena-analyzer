package usecase

import (
	"context"
	"fmt"
	"strings"

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
) *fetchUsecase {

	// target URL list
	urls := []string{
		"https://note.com/simplearchitect/n/nadc0bcdd5b3d",
		"https://note.com/simplearchitect/n/n871f29ffbfac",
		"https://note.com/simplearchitect/n/n86a95bc19b4c",
		"https://note.com/simplearchitect/n/nfd147540e3db",
	}

	return &fetchUsecase{
		logger:          logger,
		bookmarkRepo:    bookmarkRepo,
		bookmarkFetcher: bookmarkFetcher,
		urls:            urls,
	}
}

func (f *fetchUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer f.bookmarkRepo.Close()

	for _, url := range f.urls {
		// load existing data
		existingBookmark, err := f.bookmarkRepo.ReadEntity(ctx, url)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				// 初回実行時の初期化
				existingBookmark.Users = make(map[string]entities.User)
			} else {
				f.logger.Error("failed to call bookmarkRepo.ReadEntity()", "url", url, "error", err)
				continue
			}
		}

		// 既存ユーザーをすべて`isDeleted = true`に設定
		for userName := range existingBookmark.Users {
			existingBookmark.Users[userName] = entities.User{
				Name:        userName,
				IsDeleted:   true,
				IsCommented: existingBookmark.Users[userName].IsCommented,
			}
		}

		// 新しいデータを取得
		newBookmark, err := f.bookmarkFetcher.Entity(url)
		if err != nil {
			f.logger.Error("failed to call fetchBookmarkData()", "url", url, "error", err)
			continue
		}

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

		// データを保存
		err = f.bookmarkRepo.WriteEntity(ctx, url, existingBookmark)
		if err != nil {
			f.logger.Error("failed to call bookmarkRepo.WriteEntity()", "url", url, "error", err)
			continue
		}
		f.logger.Info("data saved", "url", url)

		// 表示
		fmt.Println("===================================================================")
		fmt.Printf("Title: %s\n", existingBookmark.Title)
		fmt.Printf("Count: %d\n", existingBookmark.Count)
		fmt.Printf("UserCount: %d\n", len(existingBookmark.Users))

		//fmt.Printf("Users:\n")
		// for _, user := range existingBookmark.Users {
		// 	fmt.Printf(" - %s\n", user.Name)
		// }
		fmt.Println()
	}

	return nil
}
