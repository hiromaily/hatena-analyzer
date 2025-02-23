package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/semaphore"

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
	bookmarkRepo      repository.FetchBookmarkRepositorier
	entityJSONFetcher fetcher.EntityJSONFetcher
	maxWorker         int64 // for semaphore
	urls              []string
}

// TODO
// - add cli parameter: category_code
// - fetch bookmark concurrently

func NewFetchBookmarkUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	bookmarkRepo repository.FetchBookmarkRepositorier,
	entityJSONFetcher fetcher.EntityJSONFetcher,
	maxWorker int64,
	urls []string,
) (*fetchBookmarkUsecase, error) {
	// validation
	if maxWorker == 0 {
		return nil, errors.New("maxWorker is 0")
	}

	return &fetchBookmarkUsecase{
		logger:            logger,
		tracer:            tracer,
		bookmarkRepo:      bookmarkRepo,
		entityJSONFetcher: entityJSONFetcher,
		maxWorker:         maxWorker,
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
	var entityURLs []entities.URL
	if len(f.urls) == 0 {
		var err error
		entityURLs, err = f.bookmarkRepo.GetAllURLs(ctx)
		if err != nil {
			f.logger.Error("failed to call bookmarkRepo.GetAllURLs()", "error", err)
			return err
		}
	} else {
		for _, url := range f.urls {
			entityURLs = append(entityURLs, entities.URL{Address: url})
		}
	}

	return f.concurrentExecuter(ctx, entityURLs)
}

func (f *fetchBookmarkUsecase) concurrentExecuter(
	ctx context.Context,
	entityURLs []entities.URL,
) error {
	sem := semaphore.NewWeighted(f.maxWorker)
	var wg sync.WaitGroup

	f.logger.Info("start concurrentExecuter", "max_worker", f.maxWorker, "url_count", len(entityURLs))

	for _, entityURL := range entityURLs {
		wg.Add(1)

		// get semaphore
		if err := sem.Acquire(ctx, 1); err != nil {
			f.logger.Warn("failed to acquire semaphore", "error", err)
			break
		}

		go func(entityURL entities.URL) {
			defer func() {
				wg.Done()
				sem.Release(1)
			}()

			// load existing bookmark data from DB
			existingBookmark, err := f.load(ctx, entityURL.Address)
			if err != nil {
				return
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
			newBookmark, err := f.fetch(ctx, entityURL.Address)
			if err != nil {
				return
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
				"url", entityURL.Address,
				"newBookmark.Title", existingBookmark.Title,
				"newBookmark.Count", existingBookmark.Count,
				"newBookmark.User.Length", len(existingBookmark.Users),
			)

			// save data
			err = f.save(ctx, &entityURL, existingBookmark)
			if err != nil {
				return
			}

			// Print data
			f.print(existingBookmark)
		}(entityURL)
	}
	wg.Wait()

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
	entityURL *entities.URL,
	bookmark *entities.Bookmark,
) error {
	// InfluxDB
	err := f.bookmarkRepo.WriteEntitySummary(ctx, entityURL.Address, bookmark)
	if err != nil {
		f.logger.Error(
			"failed to call bookmarkRepo.WriteEntitySummary()",
			"url", entityURL.Address,
			"error", err,
		)
		return err
	}

	// MongoDB
	err = f.bookmarkRepo.WriteEntity(ctx, entityURL.Address, bookmark)
	if err != nil {
		f.logger.Error("failed to call bookmarkRepo.WriteEntity()", "url", entityURL.Address, "error", err)
		return err
	}

	// Upsert URL to PostgreSQL DB
	if entityURL.ID == 0 { // url comes from environment variable
		urlID, err := f.bookmarkRepo.UpsertURL(
			ctx,
			entityURL.Address,
			entities.Knowledge,
			bookmark.Title,
			bookmark.Count,
			len(bookmark.Users),
			entities.PrivateUserRate(bookmark.Count, len(bookmark.Users)),
		)
		if err != nil && !rdb.IsNoRows(err) {
			f.logger.Error(
				"failed to call bookmarkRepo.UpsertURL()",
				"url", entityURL.Address,
				"error", err,
			)
			return err
		}
		if urlID == 0 {
			err := errors.New("urlID is 0")
			f.logger.Error(
				"failed to call bookmarkRepo.UpsertURL()",
				"url", entityURL.Address,
				"error", err,
			)
			return err
		}
		entityURL.ID = urlID
	} else {
		// update by urlID
		// TODO: change to update by urlAddress ??
		_, err := f.bookmarkRepo.UpdateURL(
			ctx,
			entityURL.ID,
			bookmark.Title,
			bookmark.Count,
			len(bookmark.Users),
			entities.PrivateUserRate(bookmark.Count, len(bookmark.Users)),
		)
		if err != nil && !rdb.IsNoRows(err) {
			f.logger.Error(
				"failed to call bookmarkRepo.UpdateURL()",
				"urlID", entityURL.ID,
				"error", err,
			)
			return err
		}
	}

	// Upsert Users related to url on PostgreSQL DB
	for _, users := range bookmark.Users {
		// Users
		userID, err := f.bookmarkRepo.UpsertUser(ctx, users.Name)
		if err != nil {
			f.logger.Warn("failed to call bookmarkRepo.UpsertUser()", "userName", users.Name, "error", err)
		}
		// UserURLs
		err = f.bookmarkRepo.UpsertUserURLs(ctx, userID, entityURL.ID)
		if err != nil {
			// FIXED: ERROR: insert or update on table "userurls" violates foreign key
			// constraint "userurls_url_id_fkey" (SQLSTATE 23503)
			f.logger.Warn(
				"failed to call bookmarkRepo.UpsertUserURLs()",
				"userID", userID,
				"urlID", entityURL.ID,
				"error", err,
			)
		}
	}

	f.logger.Info("bookmark data saved", "url", entityURL.Address)
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
