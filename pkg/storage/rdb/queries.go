package rdb

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/hiromaily/hatena-fake-detector/pkg/adapter"
	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb/sqlcgen"
)

type PostgreQueries struct {
	logger    logger.Logger
	rdbClient *SqlcPostgresClient
}

func NewPostgreQueries(
	logger logger.Logger,
	rdbClient *SqlcPostgresClient,
) *PostgreQueries {
	return &PostgreQueries{
		logger:    logger,
		rdbClient: rdbClient,
	}
}

func (p *PostgreQueries) Close(ctx context.Context) error {
	return p.rdbClient.Close(ctx)
}

//
// urls
//

func (p *PostgreQueries) GetURLsByURLAddresses(ctx context.Context, urls []string) ([]entities.URL, error) {
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	urlsRow, err := queries.GetURLsByURLAddresses(ctx, urls)
	if err != nil {
		return nil, err
	}
	// convert to entity models
	// []sqlcgen.GetURLsByURLAddressesRow
	urlModels := adapter.URLsByURLAddressesToEntityModel(urlsRow)
	return urlModels, nil
}

func (p *PostgreQueries) GetAllURLs(ctx context.Context) ([]entities.URLIDAddress, error) {
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	urlsRow, err := queries.GetAllURLs(ctx)
	if err != nil {
		return nil, err
	}
	// convert to entity models
	// []sqlcgen.GetAllURLsRow
	urls := adapter.URLIDAddressesToEntityModel(urlsRow)
	return urls, nil
}

func (p *PostgreQueries) GetURLID(ctx context.Context, url string) (int32, error) {
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return 0, err
	}
	defer release()
	return queries.GetUrlID(ctx, url)
}

func (p *PostgreQueries) InsertURL(
	ctx context.Context,
	url string,
	categoryCode entities.CategoryCode,
	bmCount, userCount int,
) (int32, error) {
	if url == "" || categoryCode == "" {
		return 0, errors.New("url or category code is empty")
	}

	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return 0, err
	}
	defer release()

	params := sqlcgen.InsertURLParams{
		UrlAddress:     url,
		CategoryCode:   pgtype.Text{String: categoryCode.String(), Valid: categoryCode != ""},
		BookmarkCount:  pgtype.Int4{Int32: int32(bmCount), Valid: true},
		NamedUserCount: pgtype.Int4{Int32: int32(userCount), Valid: true},
	}
	return queries.InsertURL(ctx, params)
}

func (p *PostgreQueries) UpsertURL(
	ctx context.Context,
	url string,
	categoryCode entities.CategoryCode,
	title string,
	bmCount, userCount int,
	privateUserRate float64,
) (int32, error) {
	if url == "" || categoryCode == "" {
		return 0, errors.New("url or category code is empty")
	}

	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return 0, err
	}
	defer release()

	params := sqlcgen.UpsertURLParams{
		UrlAddress:      url,
		CategoryCode:    pgtype.Text{String: categoryCode.String(), Valid: true},
		Title:           pgtype.Text{String: title, Valid: true},
		BookmarkCount:   pgtype.Int4{Int32: int32(bmCount), Valid: true},
		NamedUserCount:  pgtype.Int4{Int32: int32(userCount), Valid: true},
		PrivateUserRate: pgtype.Float8{Float64: privateUserRate, Valid: true},
	}
	return queries.UpsertURL(ctx, params)
}

func (p *PostgreQueries) UpdateURL(
	ctx context.Context,
	urlID int32,
	title string,
	bmCount, userCount int,
	privateUserRate float64,
) (int64, error) {
	if urlID == 0 {
		return 0, errors.New("urlID is 0")
	}

	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return 0, err
	}
	defer release()

	params := sqlcgen.UpdateURLParams{
		Title:           pgtype.Text{String: title, Valid: true},
		BookmarkCount:   pgtype.Int4{Int32: int32(bmCount), Valid: true},
		NamedUserCount:  pgtype.Int4{Int32: int32(userCount), Valid: true},
		PrivateUserRate: pgtype.Float8{Float64: privateUserRate, Valid: true},
		UrlID:           urlID,
	}
	return queries.UpdateURL(ctx, params)
}

func (p *PostgreQueries) InsertURLs(
	ctx context.Context,
	category entities.CategoryCode,
	urls []string,
) (int64, error) {
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return 0, err
	}
	defer release()

	params := adapter.CreateInsertURLsParams(category.String(), urls)
	return queries.InsertURLs(ctx, params)
}

//
// users
//

// func (r *PostgreQueries) InsertUser(ctx context.Context, userName string) error {
// 	return r.rdbClient.Queries.InsertUser(ctx, userName)
// }

func (p *PostgreQueries) UpsertUser(ctx context.Context, userName string) (int32, error) {
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return 0, err
	}
	defer release()
	return queries.UpsertUser(ctx, userName)
}

func (p *PostgreQueries) GetUsersByURL(ctx context.Context, url string) ([]entities.RDBUser, error) {
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	users, err := queries.GetUsersByURL(ctx, url)
	if err != nil {
		return nil, err
	}
	// convert to entity models
	return adapter.DBUsersToEntityModel(users), nil
}

func (p *PostgreQueries) GetUserNames(ctx context.Context) ([]string, error) {
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	return queries.GetUserNames(ctx)
}

func (p *PostgreQueries) GetUserNamesByURLS(ctx context.Context, urls []string) ([]string, error) {
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	return queries.GetUserNamesByURLs(ctx, urls)
}

func (p *PostgreQueries) UpdateUserBookmarkCount(ctx context.Context, userName string, count int) error {
	param := sqlcgen.UpdateUserBookmarkCountParams{
		BookmarkCount: pgtype.Int4{Int32: int32(count), Valid: true},
		UserName:      userName,
	}
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return err
	}
	defer release()
	_, err = queries.UpdateUserBookmarkCount(ctx, param)
	return err
}

//
// user_urls
//

func (p *PostgreQueries) UpsertUserURLs(ctx context.Context, userID, urlID int32) error {
	param := sqlcgen.UpsertUserURLsParams{
		UserID: userID,
		UrlID:  urlID,
	}
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return err
	}
	defer release()
	return queries.UpsertUserURLs(ctx, param)
}
