package adapter

import (
	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb/sqlcgen"
)

func DBURLsToEntityModel(urls []sqlcgen.GetAllURLsRow) []entities.RDBURL {
	var urlModels []entities.RDBURL
	for _, url := range urls {
		urlModels = append(urlModels, entities.RDBURL{
			URLID:      url.UrlID,
			URLAddress: url.UrlAddress,
		})
	}
	return urlModels
}
