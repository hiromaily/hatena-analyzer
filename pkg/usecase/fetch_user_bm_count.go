package usecase

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/semaphore"

	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
)

type FetchUserBookmarkCountUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchUserBookmarkCountUsecase struct {
	logger             logger.Logger
	tracer             tracer.Tracer
	fetchUserRepo      repository.FetchUserRepositorier
	userBMCountFetcher fetcher.UserBookmarkCountFetcher
	maxWorker          int64 // for semaphore
	urls               []string
}

func NewFetchUserBookmarkCountUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	fetchUserRepo repository.FetchUserRepositorier,
	userBMCountFetcher fetcher.UserBookmarkCountFetcher,
	maxWorker int64,
	urls []string,
) (*fetchUserBookmarkCountUsecase, error) {
	if maxWorker == 0 {
		return nil, errors.New("maxWorker is 0")
	}

	return &fetchUserBookmarkCountUsecase{
		logger:             logger,
		tracer:             tracer,
		fetchUserRepo:      fetchUserRepo,
		userBMCountFetcher: userBMCountFetcher,
		maxWorker:          maxWorker,
		urls:               urls,
	}, nil
}

// Fetch user's bookmark count of given urls by scraping
// Then save data to DB

func (f *fetchUserBookmarkCountUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer f.fetchUserRepo.Close(ctx)

	_, span := f.tracer.NewSpan(ctx, "fetchUserBookmarkCountUsecase:Execute()")
	defer func() {
		span.End()
		f.tracer.Close(ctx)
	}()

	// get user list from DB
	var users []string
	var err error
	if len(f.urls) == 0 {
		users, err = f.fetchUserRepo.GetUserNames(ctx)
		if err != nil {
			f.logger.Error("failed to get users", "error", err)
			return err
		}
	} else {
		users, err = f.fetchUserRepo.GetUserNamesByURLS(ctx, f.urls)
		if err != nil {
			f.logger.Error("failed to get users by urls", "error", err)
			return err
		}
	}
	// fetch user's bookmark count of given urls by scraping
	return f.concurrentExecuter(ctx, users)
}

func (f *fetchUserBookmarkCountUsecase) concurrentExecuter(ctx context.Context, users []string) error {
	sem := semaphore.NewWeighted(f.maxWorker)
	var wg sync.WaitGroup

	f.logger.Info("start concurrentExecuter", "max_worker", f.maxWorker, "user_count", len(users))

	for _, userName := range users {
		wg.Add(1)

		// get semaphore
		if err := sem.Acquire(ctx, 1); err != nil {
			f.logger.Warn("failed to acquire semaphore", "error", err)
			break
		}

		go func(userName string) {
			defer func() {
				wg.Done()
				sem.Release(1)
			}()

			// 1. get user's bookmark count
			bmCount, err := f.userBMCountFetcher.Fetch(ctx, userName)
			if err != nil {
				f.logger.Error("failed to get user bookmark count", "user_name", userName, "error", err)
				return
			}
			// s.logger.Debug("user info", "user_name", userName, "bm_count", bmCount)

			// 2. save data to DB
			if err := f.fetchUserRepo.UpdateUserBookmarkCount(ctx, userName, bmCount); err != nil {
				//FIXED: failed to deallocate cached statement(s): conn busy
				f.logger.Error("failed to update user bookmark count", "user_name", userName, "error", err)
				return
			}
		}(userName)
	}
	wg.Wait()

	return nil
}
