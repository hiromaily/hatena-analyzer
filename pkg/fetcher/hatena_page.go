package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
)

type hatenaPageURLFetcher struct {
	logger logger.Logger
}

func NewHatenaPageURLFetcher(logger logger.Logger) *hatenaPageURLFetcher {
	return &hatenaPageURLFetcher{
		logger: logger,
	}
}

// Fetch bookmark count of user from Hatena user's page

func (h *hatenaPageURLFetcher) Fetch(ctx context.Context, url string) ([]string, error) {
	// h.logger.Debug("hatenaPageURLFetcher.Fetch() fetching urls of page: ", "url", url)

	// Request
	resp, err := Request(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.logger.Error("failed to get page", "status_code", resp.StatusCode, "url", url)
		return nil, fmt.Errorf("failed to get user: status: %d", resp.StatusCode)
	}

	// Parse
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// h3 class=entrylist-contents-title > a href
	className := "entrylist-contents-title"
	hrefs := extractHrefsFromClass(doc, className)

	return hrefs, nil
}

func extractHrefsFromClass(node *html.Node, className string) []string {
	var hrefs []string
	var fn func(n *html.Node)
	fn = func(n *html.Node) {
		// If the node is h3 tag
		if n.Type == html.ElementNode && n.Data == "h3" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, className) {
					// Now look for 'a' tag inside this 'h3'
					var aHref func(*html.Node) string
					aHref = func(n *html.Node) string {
						if n.Type == html.ElementNode && n.Data == "a" {
							for _, attr := range n.Attr {
								if attr.Key == "href" {
									return attr.Val
								}
							}
						}
						// Recursively look into child nodes
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							if href := aHref(c); href != "" {
								return href
							}
						}
						return ""
					}
					if href := aHref(n); href != "" {
						hrefs = append(hrefs, href)
					}
				}
			}
		}

		// Continue traversing the HTML tree
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}

	fn(node)
	return hrefs
}
