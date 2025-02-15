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
)

type FetchBookmarkUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchBookmarkUsecase struct {
	logger          logger.Logger
	bookmarkRepo    repository.BookmarkRepositorier
	bookmarkFetcher fetcher.BookmarkFetcher
	urls            []string
}

func NewFetchBookmarkUsecase(
	logger logger.Logger,
	bookmarkRepo repository.BookmarkRepositorier,
	bookmarkFetcher fetcher.BookmarkFetcher,
	urls []string,
) (*fetchBookmarkUsecase, error) {
	// validation
	if len(urls) == 0 {
		return nil, errors.New("urls is empty")
	}

	return &fetchBookmarkUsecase{
		logger:          logger,
		bookmarkRepo:    bookmarkRepo,
		bookmarkFetcher: bookmarkFetcher,
		urls:            urls,
	}, nil
}

func (f *fetchBookmarkUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer f.bookmarkRepo.Close(ctx)

	for _, url := range f.urls {
		// load
		existingBookmark, err := f.load(ctx, url)
		if err != nil {
			continue
		}

		// set isDeleted = `true` on existingBookmark.Users
		for userName := range existingBookmark.Users {
			existingBookmark.Users[userName] = entities.User{
				Name:        userName,
				IsDeleted:   true,
				IsCommented: existingBookmark.Users[userName].IsCommented,
			}
		}

		// retrieve latest data from web
		newBookmark, err := f.fetch(ctx, url)
		if err != nil {
			continue
		}

		// update
		existingBookmark.Title = newBookmark.Title
		existingBookmark.Count = newBookmark.Count
		existingBookmark.Timestamp = newBookmark.Timestamp
		// overwrite `isDeleted` with `false` if user is still exist
		for userName, user := range newBookmark.Users {
			existingBookmark.Users[userName] = entities.User{
				Name:        userName,
				IsDeleted:   false,
				IsCommented: user.IsCommented,
			}
		}
		f.logger.Info("bookmark entity will be stored",
			"url", url,
			"newBookmark.Title", existingBookmark.Title,
			"newBookmark.Count", existingBookmark.Count,
			"newBookmark.User.Length", len(existingBookmark.Users),
		)

		// save data
		err = f.save(ctx, url, existingBookmark)
		if err != nil {
			continue
		}

		// Print data
		f.print(existingBookmark)
	}

	return nil
}

func (f *fetchBookmarkUsecase) load(ctx context.Context, url string) (*entities.Bookmark, error) {
	// load bookmark summary
	bookmarkSummary, err := f.bookmarkRepo.ReadEntitySummary(ctx, url)
	if err != nil {
		f.logger.Error("failed to call bookmarkRepo.ReadEntitySummary()", "url", url, "error", err)
		return nil, err
	}
	f.logger.Debug("bookmark summary loaded",
		"url", url,
		"bookmarkSummary.Title", bookmarkSummary.Title,
		"bookmarkSummary.Count", bookmarkSummary.Count,
		"bookmarkSummary.UserCount", bookmarkSummary.UserCount,
	)

	// load bookmark content by URL
	existingBookmark, err := f.bookmarkRepo.ReadEntity(ctx, url)
	if err != nil {
		f.logger.Error("failed to call bookmarkRepo.ReadEntity()", "url", url, "error", err)
		return nil, err
	}

	if existingBookmark == nil {
		f.logger.Debug("entity not found on DB", "url", url)
		// initialize entities.Bookmark
		existingBookmark = &entities.Bookmark{}
		existingBookmark.Users = make(map[string]entities.User)
	}

	f.logger.Info("bookmark entity loaded",
		"url", url,
		"existingBookmark.Title", existingBookmark.Title,
		"existingBookmark.Count", existingBookmark.Count,
		"existingBookmark.User.Length", len(existingBookmark.Users),
	)
	return existingBookmark, nil
}

func (f *fetchBookmarkUsecase) fetch(ctx context.Context, url string) (*entities.Bookmark, error) {
	newBookmark, err := f.bookmarkFetcher.Entity(ctx, url)
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

func (f *fetchBookmarkUsecase) save(ctx context.Context, url string, bookmark *entities.Bookmark) error {
	// InfluxDB
	err := f.bookmarkRepo.WriteEntitySummary(ctx, url, bookmark)
	if err != nil {
		f.logger.Error("failed to call bookmarkRepo.WriteEntitySummary()", "url", url, "error", err)
		return err
	}

	// MongoDB
	err = f.bookmarkRepo.WriteEntity(ctx, url, bookmark)
	if err != nil {
		f.logger.Error("failed to call bookmarkRepo.WriteEntity()", "url", url, "error", err)
		return err
	}

	// PostgreSQL
	// urlID, err := f.bookmarkRepo.GetURLID(ctx, url)
	// if err != nil {
	// 	f.logger.Error("failed to call bookmarkRepo.GetURLID()", "url", url, "error", err)
	// 	return err
	// }

	urlID, err := f.bookmarkRepo.InsertURL(ctx, url)
	if err != nil && !rdb.IsNoRows(err) {
		f.logger.Error("failed to call bookmarkRepo.InsertURL()", "url", url, "error", err)
		return err
	}
	if urlID == 0 {
		// TODO: urlIDを取得する必要がある
		err := errors.New("urlID is 0")
		f.logger.Error("failed to call bookmarkRepo.InsertURL()", "url", url, "error", err)
		return err
	}

	for _, users := range bookmark.Users {
		// Users
		userID, err := f.bookmarkRepo.UpsertUser(ctx, users.Name)
		if err != nil {
			f.logger.Warn("failed to call bookmarkRepo.UpsertUser()", "userName", users.Name, "error", err)
		}
		// UserURLs
		err = f.bookmarkRepo.UpsertUserURLs(ctx, userID, urlID)
		if err != nil {
			// FIXME: ERROR: insert or update on table "userurls" violates foreign key
			// constraint "userurls_url_id_fkey" (SQLSTATE 23503)
			f.logger.Warn(
				"failed to call bookmarkRepo.UpsertUserURLs()",
				"userID", userID,
				"urlID", urlID,
				"error", err,
			)
		}
	}

	f.logger.Info("bookmark data saved", "url", url)
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
