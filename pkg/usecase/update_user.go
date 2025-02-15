package usecase

import (
	"context"

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
}

func NewUpdateUserInfoUsecase(
	logger logger.Logger,
	userRepo repository.UserRepositorier,
	userFetcher fetcher.UserBookmarkFetcher,
) *updateUserInfoUsecase {
	return &updateUserInfoUsecase{
		logger:      logger,
		userRepo:    userRepo,
		userFetcher: userFetcher,
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
	for idx, userName := range users {
		// 2-1. ユーザーのブックマーク数を取得
		bmCount, err := s.userFetcher.UserBookmark(ctx, userName)
		if err != nil {
			s.logger.Error("failed to get user bookmark count", "user_name", userName, "error", err)
		} else {
			s.logger.Debug("user info", "user_name", userName, "bm_count", bmCount)
		}

		if idx > 10 {
			break
		}

		// 2-2. 取得した情報をDBに保存
	}

	return nil
}
