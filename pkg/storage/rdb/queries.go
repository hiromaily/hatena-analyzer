package rdb

import (
	"context"

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

func (p *PostgreQueries) GetAllURLs(ctx context.Context) ([]entities.RDBURL, error) {
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
	urls := adapter.DBURLsToEntityModel(urlsRow)
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

func (p *PostgreQueries) InsertURL(ctx context.Context, url string) (int32, error) {
	queries, release, err := p.rdbClient.GetQueries(ctx)
	if err != nil {
		return 0, err
	}
	defer release()
	return queries.InsertURL(ctx, url)
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
