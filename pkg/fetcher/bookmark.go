package fetcher

import (
	"context"
	"encoding/json"
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

// FIXME: integrate bookmark entity
type Bookmark struct {
	User      string `json:"user"`
	Comment   string `json:"comment"`
	Timestamp string `json:"timestamp"`
}

type Data struct {
	Title     string     `json:"title"`
	Count     int        `json:"count"`
	Bookmarks []Bookmark `json:"bookmarks"`
}

func (b *bookmarkFetcher) Entity(ctx context.Context, url string) (*entities.Bookmark, error) {
	apiURL := b.entityURL + url

	// warning: net/http.Get must not be called (noctx)
	// resp, err := http.Get(apiURL)
	resp, err := Request(ctx, apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data Data
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
