package fetcher

import (
	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
)

type BookmarkFetcher interface {
	Entity(url string) (*entities.Bookmark, error)
}
