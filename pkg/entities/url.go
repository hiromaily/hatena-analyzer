package entities

type RDBURL struct {
	URLID      int32
	URLAddress string
}

func FilterURLAddress(urls []RDBURL) []string {
	urlAddresses := make([]string, 0, len(urls))
	for _, url := range urls {
		urlAddresses = append(urlAddresses, url.URLAddress)
	}
	return urlAddresses
}
