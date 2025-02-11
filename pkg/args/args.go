package args

import (
	"fmt"

	"github.com/alexflint/go-arg"

	"github.com/hiromaily/hatena-fake-detector/pkg/app"
)

type SubCommand struct{}

type Args struct {
	Version bool     // global option
	URLs    []string `arg:"--urls,env:URLS"` // global option

	FetchCommand         *SubCommand `arg:"subcommand:fetch-bookmark"`  // fetch latest bookmark entity
	ViewSummaryCommand   *SubCommand `arg:"subcommand:view-summary"`    // view time series bookmark summary
	FetchUserInfoCommand *SubCommand `arg:"subcommand:fetch-user-info"` // fetch user info from bookmark url
}

func Parse() (*Args, *arg.Parser, app.AppCode) {
	var args Args
	p := arg.MustParse(&args)
	return &args, p, getAppCode(&args)
}

func getAppCode(args *Args) app.AppCode {
	switch {
	case args.FetchCommand != nil:
		return app.AppCodeFetchBookmark
	case args.ViewSummaryCommand != nil:
		return app.AppCodeViewSummary
	case args.FetchUserInfoCommand != nil:
		return app.AppCodeFetchUserInfo
	}
	panic(fmt.Errorf("subcommand is not found"))
}
