package app

type AppCode string

// 適宜コマンドの実装に合わせて追加
const (
	AppCodeFetch         = AppCode("Fetch")
	AppCodeViewSummary   = AppCode("ViewSummary")
	AppCodePrintBookmark = AppCode("PrintBookmark")
	AppCodeWeb           = AppCode("WebServer")
)

func (a AppCode) String() string {
	return string(a)
}
