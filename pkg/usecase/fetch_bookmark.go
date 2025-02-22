package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
)

type FetchBookmarkUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchBookmarkUsecase struct {
	logger            logger.Logger
	tracer            tracer.Tracer
	bookmarkRepo      repository.BookmarkRepositorier
	entityJSONFetcher fetcher.EntityJSONFetcher
	urls              []string
}

// TODO
// - add cli parameter: category_code

func NewFetchBookmarkUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	bookmarkRepo repository.BookmarkRepositorier,
	entityJSONFetcher fetcher.EntityJSONFetcher,
	urls []string,
) (*fetchBookmarkUsecase, error) {
	// validation

	return &fetchBookmarkUsecase{
		logger:            logger,
		tracer:            tracer,
		bookmarkRepo:      bookmarkRepo,
		entityJSONFetcher: entityJSONFetcher,
		urls:              urls,
	}, nil
}

// Fetch bookmark users, title, count related given URLs using Hatena entity API and save data to DB

func (f *fetchBookmarkUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer f.bookmarkRepo.Close(ctx)

	_, span := f.tracer.NewSpan(ctx, "fetchBookmarkUsecase:Execute()")
	defer func() {
		span.End()
		f.tracer.Close(ctx)
	}()

	// get urls from DB if needed
	var entityURLs []entities.RDBURL
	if len(f.urls) == 0 {
		var err error
		entityURLs, err = f.bookmarkRepo.GetAllURLs(ctx)
		if err != nil {
			f.logger.Error("failed to call bookmarkRepo.GetAllURLs()", "error", err)
			return err
		}
		// f.urls = entities.FilterURLAddress(entityURLs)
		// isDBURLs = true
	} else {
		for _, url := range f.urls {
			entityURLs = append(entityURLs, entities.RDBURL{URLAddress: url})
		}
	}

	for _, entityURL := range entityURLs {
		// load existing bookmark data from DB
		existingBookmark, err := f.load(ctx, entityURL.URLAddress)
		if err != nil {
			continue
		}

		// set isDeleted = `true` on existingBookmark.Users
		for userName := range existingBookmark.Users {
			existingBookmark.Users[userName] = entities.BookmarkUser{
				Name:        userName,
				IsDeleted:   true,
				IsCommented: existingBookmark.Users[userName].IsCommented,
			}
		}

		// retrieve latest data from URL
		newBookmark, err := f.fetch(ctx, entityURL.URLAddress)
		if err != nil {
			continue
		}

		// update existingBookmark to save
		existingBookmark.Title = newBookmark.Title
		existingBookmark.Count = newBookmark.Count
		existingBookmark.Timestamp = newBookmark.Timestamp
		// overwrite `isDeleted` with `false` if user is still exist
		for userName, user := range newBookmark.Users {
			existingBookmark.Users[userName] = entities.BookmarkUser{
				Name:        userName,
				IsDeleted:   false,
				IsCommented: user.IsCommented,
			}
		}
		f.logger.Info("bookmark entity will be stored",
			"url", entityURL.URLAddress,
			"newBookmark.Title", existingBookmark.Title,
			"newBookmark.Count", existingBookmark.Count,
			"newBookmark.User.Length", len(existingBookmark.Users),
		)

		// save data
		err = f.save(ctx, &entityURL, existingBookmark)
		if err != nil {
			continue
		}

		// Print data
		f.print(existingBookmark)
	}

	return nil
}

// load existing bookmark data from DB
func (f *fetchBookmarkUsecase) load(ctx context.Context, url string) (*entities.Bookmark, error) {
	// load bookmark summary from InfluxDB
	// bookmarkSummary, err := f.bookmarkRepo.ReadEntitySummary(ctx, url)
	// if err != nil {
	// 	f.logger.Error("failed to call bookmarkRepo.ReadEntitySummary()", "url", url, "error", err)
	// 	return nil, err
	// }
	// f.logger.Debug("bookmark summary loaded",
	// 	"url", url,
	// 	"bookmarkSummary.Title", bookmarkSummary.Title,
	// 	"bookmarkSummary.Count", bookmarkSummary.Count,
	// 	"bookmarkSummary.UserCount", bookmarkSummary.UserCount,
	// )

	// load bookmark content by URL from MongoDB
	existingBookmark, err := f.bookmarkRepo.ReadEntity(ctx, url)
	if err != nil {
		f.logger.Error("failed to call bookmarkRepo.ReadEntity()", "url", url, "error", err)
		return nil, err
	}

	if existingBookmark == nil {
		f.logger.Debug("entity not found on MongoDB", "url", url)
		// initialize entities.Bookmark
		existingBookmark = &entities.Bookmark{}
		existingBookmark.Users = make(map[string]entities.BookmarkUser)
	} else {
		f.logger.Info("bookmark entity loaded",
			"url", url,
			"existingBookmark.Title", existingBookmark.Title,
			"existingBookmark.Count", existingBookmark.Count,
			"existingBookmark.User.Length", len(existingBookmark.Users),
		)
	}

	return existingBookmark, nil
}

func (f *fetchBookmarkUsecase) fetch(ctx context.Context, url string) (*entities.Bookmark, error) {
	newBookmark, err := f.entityJSONFetcher.Fetch(ctx, url)
	if err != nil {
		f.logger.Error("failed to call fetchBookmarkData()", "url", url, "error", err)
		return nil, err
	}
	f.logger.Info(
		"data fetched",
		"url", url,
		"newBookmark.Title", newBookmark.Title,
		"newBookmark.Count", newBookmark.Count,
		"newBookmark.User.Length", len(newBookmark.Users),
	)

	return newBookmark, nil
}

func (f *fetchBookmarkUsecase) save(
	ctx context.Context,
	entityURL *entities.RDBURL,
	bookmark *entities.Bookmark,
) error {
	// InfluxDB
	err := f.bookmarkRepo.WriteEntitySummary(ctx, entityURL.URLAddress, bookmark)
	if err != nil {
		f.logger.Error(
			"failed to call bookmarkRepo.WriteEntitySummary()",
			"url", entityURL.URLAddress,
			"error", err,
		)
		return err
	}

	// MongoDB
	err = f.bookmarkRepo.WriteEntity(ctx, entityURL.URLAddress, bookmark)
	if err != nil {
		f.logger.Error("failed to call bookmarkRepo.WriteEntity()", "url", entityURL.URLAddress, "error", err)
		return err
	}

	// Insert URL to PostgreSQL DB
	if entityURL.URLID == 0 { // url comes from environment variable
		urlID, err := f.bookmarkRepo.InsertURL(ctx, entityURL.URLAddress)
		if err != nil && !rdb.IsNoRows(err) {
			f.logger.Error(
				"failed to call bookmarkRepo.InsertURL()",
				"url", entityURL.URLAddress,
				"error", err,
			)
			return err
		}
		if urlID == 0 {
			err := errors.New("urlID is 0")
			f.logger.Error(
				"failed to call bookmarkRepo.InsertURL()",
				"url", entityURL.URLAddress,
				"error", err,
			)
			return err
		}
		entityURL.URLID = urlID
	}

	// Upsert Users related to url on PostgreSQL DB
	for _, users := range bookmark.Users {
		// Users
		userID, err := f.bookmarkRepo.UpsertUser(ctx, users.Name)
		if err != nil {
			f.logger.Warn("failed to call bookmarkRepo.UpsertUser()", "userName", users.Name, "error", err)
		}
		// UserURLs
		err = f.bookmarkRepo.UpsertUserURLs(ctx, userID, entityURL.URLID)
		if err != nil {
			// FIXED: ERROR: insert or update on table "userurls" violates foreign key
			// constraint "userurls_url_id_fkey" (SQLSTATE 23503)
			f.logger.Warn(
				"failed to call bookmarkRepo.UpsertUserURLs()",
				"userID", userID,
				"urlID", entityURL.URLID,
				"error", err,
			)
		}
	}

	f.logger.Info("bookmark data saved", "url", entityURL.URLAddress)
	return nil
}

func (f *fetchBookmarkUsecase) print(bookmark *entities.Bookmark) {
	fmt.Println("===================================================================")
	fmt.Printf("Title: %s\n", bookmark.Title)
	fmt.Printf("Count: %d\n", bookmark.Count)
	fmt.Printf("UserCount: %d\n", len(bookmark.Users))
	fmt.Printf("DeletedUserCount: %d\n", bookmark.CountDeletedUser())

	// fmt.Printf("Users:\n")
	// for _, user := range existingBookmark.Users {
	// 	fmt.Printf(" - %s\n", user.Name)
	// }
	fmt.Println()
}
