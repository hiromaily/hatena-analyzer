package app

type AppCode string

// AppCode
const (
	AppCodeFetchPageURLs          = AppCode("FetchPageURLs")
	AppCodeFetchBookmarkEntities  = AppCode("FetchBookmarkEntities")
	AppCodeFetchUserBookmarkCount = AppCode("FetchUserBookmarkCount")
	AppCodeViewSummary            = AppCode("ViewSummary")

	AppCodeWeb = AppCode("WebServer")
)

func (a AppCode) String() string {
	return string(a)
}
