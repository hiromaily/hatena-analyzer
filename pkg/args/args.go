package args

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/hiromaily/hatena-fake-detector/pkg/app"
)

type SubCommand struct{}

// 適宜コマンドの実装に合わせて追加
type Args struct {
	Version                bool
	FetchCommand           *SubCommand `arg:"subcommand:fetch"`             // 最新のブックマークデータを取得
	PrintBookmarkCommand   *SubCommand `arg:"subcommand:print-bookmark"`    // 最新のブックマークデータを表示
	PrintTimeSeriesCommand *SubCommand `arg:"subcommand:print-time-series"` // 時系列ブックマークのサマリーを表示
}

func Parse() (*Args, *arg.Parser, app.AppCode) {
	var args Args
	p := arg.MustParse(&args)
	return &args, p, getAppCode(&args)
}

// 適宜コマンドの実装に合わせて追加
// Note: 利用できないCommandはこちらには追加しない
func getAppCode(args *Args) app.AppCode {
	switch {
	case args.FetchCommand != nil:
		return app.AppCodeFetch
	case args.PrintBookmarkCommand != nil:
		return app.AppCodePrintBookmark
	case args.PrintTimeSeriesCommand != nil:
		return app.AppCodePrintTimeSeries
	}
	panic(fmt.Errorf("subcommand is not found"))
}
