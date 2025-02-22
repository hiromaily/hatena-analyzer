package usecase

import (
	"context"
	"fmt"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
)

type ViewSummaryUsecaser interface {
	Execute(ctx context.Context) error
}

type summaryUsecase struct {
	logger      logger.Logger
	tracer      tracer.Tracer
	summaryRepo repository.SummaryRepositorier
	urls        []string
}

func NewViewSummaryUsecase(
	logger logger.Logger,
	tracer tracer.Tracer,
	summaryRepo repository.SummaryRepositorier,
	urls []string,
) (*summaryUsecase, error) {
	// validation
	// if len(urls) == 0 {
	// 	return nil, errors.New("urls is empty")
	// }

	return &summaryUsecase{
		logger:      logger,
		tracer:      tracer,
		summaryRepo: summaryRepo,
		urls:        urls,
	}, nil
}

func (s *summaryUsecase) Execute(ctx context.Context) error {
	// must be closed dbClient
	defer s.summaryRepo.Close(ctx)

	_, span := s.tracer.NewSpan(ctx, "summaryUsecase:Execute()")
	defer func() {
		span.End()
		s.tracer.Close(ctx)
	}()

	// get urls from DB if needed
	var entityURLs []entities.RDBURL
	if len(s.urls) == 0 {
		var err error
		entityURLs, err = s.summaryRepo.GetAllURLs(ctx)
		if err != nil {
			s.logger.Error("failed to call bookmarkRepo.GetAllURLs()", "error", err)
			return err
		}
		s.urls = entities.FilterURLAddress(entityURLs)
	}

	validURLCount := 0
	privateUserRateSum := 0.0
	for _, url := range s.urls {
		// get summaries from InfluxDB
		summaries, err := s.summaryRepo.ReadEntitySummaries(ctx, url)
		if err != nil {
			s.logger.Error("failed to call summaryRepo.ReadEntitySummaries()", "url", url, "error", err)
			continue
		}
		if len(summaries) == 0 {
			s.logger.Warn("no data", "url", url)
			continue
		}

		validURLCount++
		privateUserRate := entities.PrivateUserRate(summaries[0].Count, summaries[0].UserCount)
		privateUserRateSum += privateUserRate

		fmt.Println("-------------------------------------------------------------")
		fmt.Printf(" Title: %s,\n URL: %s\n", summaries[0].Title, url)
		fmt.Printf(" Total Bookmark: %d\n", summaries[0].Count)
		fmt.Printf(" Private User Rate: %.1f\n", privateUserRate)
		// fmt.Printf(" Time series\n")
		// for _, summary := range summaries {
		// 	fmt.Printf(
		// 		"  - %s: count: %d, user_count: %d, deleted_user_count: %d, private user rate: %.1f\n",
		// 		times.ToJPTime(summary.Timestamp).Format(time.RFC3339),
		// 		summary.Count,
		// 		summary.UserCount,
		// 		summary.DeletedUserCount,
		// 		float64(summary.Count-summary.UserCount)/float64(summary.Count)*100,
		// 	)
		// }

		// get user by URL info from DB
		// users, err := s.summaryRepo.GetUsersByURL(ctx, url)
		// if err != nil {
		// 	s.logger.Error("failed to call summaryRepo.GetUsersByURL()", "url", url, "error", err)
		// 	continue
		// }
		// count users
		// var count10, count100, count1000, count10000, countOver int
		// for _, user := range users {
		// 	switch {
		// 	case user.BookmarkCount < 10:
		// 		count10++
		// 	case user.BookmarkCount < 100:
		// 		count100++
		// 	case user.BookmarkCount < 1000:
		// 		count1000++
		// 	case user.BookmarkCount < 10000:
		// 		count10000++
		// 	default:
		// 		countOver++
		// 	}
		// }
		// // calculate average
		// // less 10 user must be suspicious
		// newUserRate := float64(count10) / float64(summaries[0].UserCount) * 100

		// fmt.Printf(" User count by bookmark count\n")
		// fmt.Printf("  less 10:      %5d\n", count10)
		// fmt.Printf("  less 100:     %5d\n", count100)
		// fmt.Printf("  less 1000:    %5d\n", count1000)
		// fmt.Printf("  less 10000:   %5d\n", count10000)
		// fmt.Printf("  over 10000:   %5d\n", countOver)
		// fmt.Printf(" New user rate:  %.1f\n", newUserRate)
	}

	// calculate average
	averagePrivateUserRate := privateUserRateSum / float64(validURLCount)
	fmt.Println("-------------------------------------------------------------")
	fmt.Printf("Average private user rate: %.1f\n", averagePrivateUserRate)

	return nil
}
