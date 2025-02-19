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

type UpdateUserInfoUsecaser interface {
	Execute(ctx context.Context) error
}

type updateUserInfoUsecase struct {
	logger      logger.Logger
	tracer      tracer.Tracer
	userRepo    repository.UserRepositorier
	userFetcher fetcher.UserBookmarkFetcher
	maxWorker   int64 // for semaphore
	urls        []string
}

func NewUpdateUserInfoUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	userRepo repository.UserRepositorier,
	userFetcher fetcher.UserBookmarkFetcher,
	maxWorker int64,
	urls []string,
) (*updateUserInfoUsecase, error) {
	if maxWorker == 0 {
		return nil, errors.New("maxWorker is 0")
	}

	return &updateUserInfoUsecase{
		logger:      logger,
		tracer:      tracer,
		userRepo:    userRepo,
		userFetcher: userFetcher,
		maxWorker:   maxWorker,
		urls:        urls,
	}, nil
}

func (s *updateUserInfoUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer s.userRepo.Close(ctx)

	_, span := s.tracer.NewSpan(ctx, "updateUserInfoUsecase:Execute()")
	defer func() {
		span.End()
		s.tracer.Close(ctx)
	}()

	// 1. DBからuser一覧を取得
	var users []string
	var err error
	if len(s.urls) == 0 {
		users, err = s.userRepo.GetUserNames(ctx)
		if err != nil {
			s.logger.Error("failed to get users", "error", err)
			return err
		}
	} else {
		users, err = s.userRepo.GetUserNamesByURLS(ctx, s.urls)
		if err != nil {
			s.logger.Error("failed to get users by urls", "error", err)
			return err
		}
	}
	// 2. 取得したuser情報からscrapingでユーザーの情報を取得してDBに保存
	return s.concurrentExecuter(ctx, users)
}

func (s *updateUserInfoUsecase) concurrentExecuter(ctx context.Context, users []string) error {
	sem := semaphore.NewWeighted(s.maxWorker)
	var wg sync.WaitGroup

	s.logger.Info("start concurrentExecuter", "max_worker", s.maxWorker, "user_count", len(users))

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
			// s.logger.Debug("user info", "user_name", userName, "bm_count", bmCount)

			// 2-2. 取得した情報をDBに保存
			if err := s.userRepo.UpdateUserBookmarkCount(ctx, userName, bmCount); err != nil {
				//FIXED: failed to deallocate cached statement(s): conn busy
				s.logger.Error("failed to update user bookmark count", "user_name", userName, "error", err)
				return
			}
		}(userName)
	}
	return nil
}
