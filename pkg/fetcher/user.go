package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
)

type bookmarkUserFetcher struct {
	logger  logger.Logger
	userURL string
}

func NewBookmarkUserFetcher(logger logger.Logger) *bookmarkUserFetcher {
	return &bookmarkUserFetcher{
		logger:  logger,
		userURL: "https://b.hatena.ne.jp/%s/",
	}
}

func (b *bookmarkUserFetcher) UserBookmark(ctx context.Context, userName string) (int, error) {
	// Request
	userURL := fmt.Sprintf(b.userURL, userName)
	resp, err := Request(ctx, userURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b.logger.Error("failed to get user", "status_code", resp.StatusCode)
		return 0, fmt.Errorf("failed to get user: status: %d", resp.StatusCode)
	}

	// Parse
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to parse HTML: %v", err)
	}

	targetClassName := "userprofile-status-count"
	count, found := extractBookmarkCount(doc, targetClassName)
	if !found {
		return 0, fmt.Errorf("failed to get user bookmark count. user: %s", userName)
	}

	return count, nil
}

func hasClass(n *html.Node, class string) bool {
	for _, a := range n.Attr {
		if a.Key == "class" {
			classes := strings.Fields(a.Val)
			for _, c := range classes {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

func extractBookmarkCount(n *html.Node, class string) (int, bool) {
	if n.Type == html.ElementNode && n.Data == "span" && hasClass(n, class) && n.FirstChild != nil {
		// remove comma first
		data := strings.Replace(n.FirstChild.Data, ",", "", -1)
		if value, err := strconv.Atoi(data); err == nil {
			return value, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if value, found := extractBookmarkCount(c, class); found {
			return value, true
		}
	}
	return 0, false
}
