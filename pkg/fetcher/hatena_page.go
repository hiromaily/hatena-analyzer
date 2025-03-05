package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
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

func (h *hatenaPageURLFetcher) Fetch(
	ctx context.Context,
	url string,
	isAll bool,
) ([]entities.LinkInfo, error) {
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
	hrefs := h.extractLinkInfoFromClass(doc, className, isAll)

	return hrefs, nil
}

func (h *hatenaPageURLFetcher) extractLinkInfoFromClass(
	node *html.Node,
	className string,
	isAll bool,
) []entities.LinkInfo {
	var linkInfos []entities.LinkInfo
	var fn func(n *html.Node)
	fn = func(n *html.Node) {
		// If the node is h3 tag
		if n.Type == html.ElementNode && n.Data == "h3" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, className) {
					// Now look for 'a' tag inside this 'h3'
					var aHref func(*html.Node) (string, string)
					aHref = func(n *html.Node) (string, string) {
						if n.Type == html.ElementNode && n.Data == "a" {
							var href, entryCategory string
							for _, attr := range n.Attr {
								if attr.Key == "href" {
									href = attr.Val
								}
								if attr.Key == "data-entry-category" {
									entryCategory = attr.Val
								}
							}
							if href != "" && entryCategory != "" {
								return href, entryCategory
							}
						}
						// Recursively look into child nodes
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							if href, entryCategory := aHref(c); href != "" && entryCategory != "" {
								return href, entryCategory
							}
						}
						return "", ""
					}
					if href, entryCategory := aHref(n); href != "" && entryCategory != "" {
						cateCode := entities.GetCategoryCode(entryCategory)
						linkInfos = append(
							linkInfos,
							entities.LinkInfo{Href: href, Category: cateCode, IsAll: isAll},
						)
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
	return linkInfos
}
