package fetcher

import (
	"context"
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
	resp, err := b.request(ctx, apiURL)
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

func (b *bookmarkFetcher) request(ctx context.Context, apiURL string) (*http.Response, error) {
	// set 10 seconds timeout
	// use new context due to multiple calls
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
