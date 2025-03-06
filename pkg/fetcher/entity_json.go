package fetcher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pingcap/errors"

	"github.com/hiromaily/hatena-analyzer/pkg/entities"
	"github.com/hiromaily/hatena-analyzer/pkg/logger"
)

type entityJSONFetcher struct {
	logger    logger.Logger
	entityURL string
}

func NewEntityJSONFetcher(logger logger.Logger) *entityJSONFetcher {
	return &entityJSONFetcher{
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

// Call Hatena entity API to get bookmark users, title, count

func (e *entityJSONFetcher) Fetch(ctx context.Context, url string) (*entities.Bookmark, error) {
	apiURL := e.entityURL + url

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

	// validation
	if data.Title == "" {
		return nil, errors.New("entity is not found")
	}

	users := make(map[string]entities.BookmarkUser)
	for _, bookmark := range data.Bookmarks {
		users[bookmark.User] = entities.BookmarkUser{
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
