package app

type AppCode string

// AppCode
const (
	AppCodeFetchHatenaPageURLs    = AppCode("FetchHatenaPageURLs")
	AppCodeFetchBookmarkEntities  = AppCode("FetchBookmarkEntities")
	AppCodeFetchUserBookmarkCount = AppCode("FetchUserBookmarkCount")
	AppCodeViewSummary            = AppCode("ViewSummary")

	AppCodeWeb = AppCode("WebServer")
)

func (a AppCode) String() string {
	return string(a)
}
