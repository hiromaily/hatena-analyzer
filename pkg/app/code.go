package app

type AppCode string

// AppCode
const (
	AppCodeFetchBookmark  = AppCode("FetchBookmark")
	AppCodeViewSummary    = AppCode("ViewSummary")
	AppCodeUpdateUserInfo = AppCode("UpdateUserInfo")
	AppCodeWeb            = AppCode("WebServer")
)

func (a AppCode) String() string {
	return string(a)
}
