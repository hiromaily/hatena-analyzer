package usecase

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"

	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
)

type UpdateUserInfoUsecaser interface {
	Execute(ctx context.Context) error
}

type updateUserInfoUsecase struct {
	logger      logger.Logger
	userRepo    repository.UserRepositorier
	userFetcher fetcher.UserBookmarkFetcher
	maxWorker   int64 // for semaphore
}

func NewUpdateUserInfoUsecase(
	logger logger.Logger,
	userRepo repository.UserRepositorier,
	userFetcher fetcher.UserBookmarkFetcher,
	maxWorker int64,
) *updateUserInfoUsecase {
	return &updateUserInfoUsecase{
		logger:      logger,
		userRepo:    userRepo,
		userFetcher: userFetcher,
		maxWorker:   maxWorker,
	}
}

func (s *updateUserInfoUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer s.userRepo.Close(ctx)

	// 1. DBからuser一覧を取得
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		s.logger.Error("failed to get users", "error", err)
		return err
	}
	// 2. 取得したuser情報からscrapingでユーザーの情報を取得してDBに保存
	return s.concurrentExecuter(ctx, users)
}

func (s *updateUserInfoUsecase) concurrentExecuter(ctx context.Context, users []string) error {
	sem := semaphore.NewWeighted(s.maxWorker)
	var wg sync.WaitGroup

	for _, userName := range users {
		wg.Add(1)

		// get semaphore
		if err := sem.Acquire(ctx, 1); err != nil {
			s.logger.Warn("failed to acquire semaphore", "error", err)
			break
		}

		go func(userName string) {
			defer func() {
				wg.Done()
				sem.Release(1)
			}()

			// 2-1. ユーザーのブックマーク数を取得
			bmCount, err := s.userFetcher.UserBookmark(ctx, userName)
			if err != nil {
				s.logger.Error("failed to get user bookmark count", "user_name", userName, "error", err)
				return
			}
			s.logger.Debug("user info", "user_name", userName, "bm_count", bmCount)

			// 2-2. 取得した情報をDBに保存
			if err := s.userRepo.UpdateUserBookmarkCount(ctx, userName, bmCount); err != nil {
				//FIXME: failed to deallocate cached statement(s): conn busy
				s.logger.Error("failed to update user bookmark count", "user_name", userName, "error", err)
				return
			}
		}(userName)
	}
	return nil
}
