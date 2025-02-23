package adapter

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb/sqlcgen"
)

func AllURLsToEntityModel(urls []sqlcgen.GetAllURLsRow) []entities.URL {
	var urlModels []entities.URL
	for _, url := range urls {
		urlModels = append(urlModels, entities.URL{
			ID:              url.UrlID,
			Address:         url.UrlAddress,
			CategoryCode:    entities.CategoryCode(url.CategoryCode.String),
			Title:           url.Title.String,
			BookmarkCount:   url.BookmarkCount.Int32,
			NamedUserCount:  url.NamedUserCount.Int32,
			PrivateUserRate: url.PrivateUserRate.Float64,
		})
	}
	return urlModels
}

func URLsByURLAddressesToEntityModel(urls []sqlcgen.GetURLsByURLAddressesRow) []entities.URL {
	var urlModels []entities.URL
	for _, url := range urls {
		urlModels = append(urlModels, entities.URL{
			ID:              url.UrlID,
			Address:         url.UrlAddress,
			CategoryCode:    entities.CategoryCode(url.CategoryCode.String),
			Title:           url.Title.String,
			BookmarkCount:   url.BookmarkCount.Int32,
			NamedUserCount:  url.NamedUserCount.Int32,
			PrivateUserRate: url.PrivateUserRate.Float64,
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
