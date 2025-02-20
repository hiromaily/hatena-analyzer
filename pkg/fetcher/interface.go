package fetcher

import (
	"context"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
)

type HatenaPageURLFetcher interface {
	Fetch(ctx context.Context, url string) ([]string, error)
}

type EntityJSONFetcher interface {
	Fetch(ctx context.Context, url string) (*entities.Bookmark, error)
}

type UserBookmarkCountFetcher interface {
	Fetch(ctx context.Context, userName string) (int, error)
}
