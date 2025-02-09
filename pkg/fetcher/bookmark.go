package fetcher

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
)

type bookmarkFetcher struct {
	logger    logger.Logger
	entityURL string
}

func NewBookmarkFetcher(logger logger.Logger) *bookmarkFetcher {
	return &bookmarkFetcher{
		logger:    logger,
		entityURL: "https://b.hatena.ne.jp/entry/json/",
	}
}

func (b *bookmarkFetcher) Entity(url string) (*entities.Bookmark, error) {
	apiUrl := b.entityURL + url

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	users := make(map[string]entities.User)
	for _, bookmark := range data.Bookmarks {
		users[bookmark.User] = entities.User{
			Name:        bookmark.User,
			IsDeleted:   false,
			IsCommented: bookmark.Comment != "",
		}
	}

	return &entities.Bookmark{
		Title:     data.Title,
		Count:     data.Count,
		Users:     users,
		Timestamp: time.Now(),
	}, nil
}
