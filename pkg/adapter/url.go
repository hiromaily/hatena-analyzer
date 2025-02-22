package adapter

import (
	"github.com/jackc/pgx/v5/pgtype"

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

func CreateInsertURLsParams(category string, urls []string) []sqlcgen.InsertURLsParams {
	var params []sqlcgen.InsertURLsParams
	pgtextCategory := pgtype.Text{
		String: category,
		Valid:  category != "",
	}
	for _, url := range urls {
		params = append(params, sqlcgen.InsertURLsParams{
			UrlAddress:   url,
			CategoryCode: pgtextCategory,
		})
	}
	return params
}
