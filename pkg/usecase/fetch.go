package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
)

type FetchUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchUsecase struct {
	logger       logger.Logger
	bookmarkRepo repository.BookmarkRepositorier
	urls         []string
}

func NewFetchUsecase(
	logger logger.Logger,
	bookmarkRepo repository.BookmarkRepositorier,
) *fetchUsecase {

	// target URL list
	urls := []string{
		"https://note.com/simplearchitect/n/nadc0bcdd5b3d",
		"https://note.com/simplearchitect/n/n871f29ffbfac",
		"https://note.com/simplearchitect/n/n86a95bc19b4c",
		"https://note.com/simplearchitect/n/nfd147540e3db",
	}

	return &fetchUsecase{
		logger:       logger,
		bookmarkRepo: bookmarkRepo,
		urls:         urls,
	}
}

func fetchBookmarkData(url string) (entities.Bookmark, error) {
	apiUrl := "https://b.hatena.ne.jp/entry/json/" + url

	resp, err := http.Get(apiUrl)
	if err != nil {
		return entities.Bookmark{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Title     string `json:"title"`
		Count     int    `json:"count"`
		Bookmarks []struct {
			User      string `json:"user"`
			Comment   string `json:"comment"`
			Timestamp string `json:"timestamp"`
		} `json:"bookmarks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return entities.Bookmark{}, err
	}

	users := make(map[string]entities.User)
	for _, bookmark := range data.Bookmarks {
		users[bookmark.User] = entities.User{
			Name:        bookmark.User,
			IsDeleted:   false,
			IsCommented: bookmark.Comment != "",
		}
	}

	return entities.Bookmark{
		Title:     data.Title,
		Count:     data.Count,
		Users:     users,
		Timestamp: time.Now(),
	}, nil
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
				fmt.Printf("Error loading existing data for %s: %v\n", url, err)
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
		newBookmark, err := fetchBookmarkData(url)
		if err != nil {
			fmt.Printf("Error fetching data for %s: %v\n", url, err)
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
			fmt.Printf("Error writing data to InfluxDB for %s: %v\n", url, err)
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
