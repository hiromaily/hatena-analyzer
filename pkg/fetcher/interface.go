package fetcher

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
)

type BookmarkFetcher interface {
	Entity(ctx context.Context, url string) (*entities.Bookmark, error)
}
