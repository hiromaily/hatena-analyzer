package args

import (
	"fmt"

	"github.com/alexflint/go-arg"

	"github.com/hiromaily/hatena-fake-detector/pkg/app"
)

type SubCommand struct{}

// type SubCmdFetchUserBookmarkCount struct {
// 	URL string `arg:"--url"`
// }

type Args struct {
	Version bool     // global option
	URLs    []string `arg:"--urls,env:URLS"` // global option

	// fetch URLs from hatena pages
	FetchHatenaPageURLsCommand *SubCommand `arg:"subcommand:fetch-hatena-page-urls"`
	// fetch bookmark entity from bookmark url
	FetchBookmarkEntitiesCommand *SubCommand `arg:"subcommand:fetch-bookmark"`
	// fetch user bookmark count from bookmark url
	FetchUserBookmarkCountCommand *SubCommand `arg:"subcommand:fetch-user-bm-count"`
	// view time series of bookmark
	ViewTimeSeriesCommand *SubCommand `arg:"subcommand:view-time-series"`
	// view bookmark details
	ViewBookmarkDetailsCommand *SubCommand `arg:"subcommand:view-bookmark-details"`
	// view bookmark summary
	ViewSummaryCommand *SubCommand `arg:"subcommand:view-summary"`
}

func Parse() (*Args, *arg.Parser, app.AppCode) {
	var args Args
	p := arg.MustParse(&args)
	return &args, p, getAppCode(&args)
}

func getAppCode(args *Args) app.AppCode {
	switch {
	case args.FetchHatenaPageURLsCommand != nil:
		return app.AppCodeFetchHatenaPageURLs
	case args.FetchBookmarkEntitiesCommand != nil:
		return app.AppCodeFetchBookmarkEntities
	case args.FetchUserBookmarkCountCommand != nil:
		return app.AppCodeFetchUserBookmarkCount
	case args.ViewTimeSeriesCommand != nil:
		return app.AppCodeViewTimeSeries
	case args.ViewBookmarkDetailsCommand != nil:
		return app.AppCodeViewBookmarkDetails
	case args.ViewSummaryCommand != nil:
		return app.AppCodeViewSummary
	}
	panic(fmt.Errorf("subcommand is not found"))
}
