package app

type AppCode string

// 適宜コマンドの実装に合わせて追加
const (
	AppCodeFetch           = AppCode("Fetch")
	AppCodePrintBookmark   = AppCode("PrintBookmark")
	AppCodePrintTimeSeries = AppCode("PrintTimeSeries")
	AppCodeWeb             = AppCode("WebServer")
)

func (a AppCode) String() string {
	return string(a)
}
