package entities

type URLIDAddress struct {
	ID      int32
	Address string
}

func FilterURLAddress(urls []URLIDAddress) []string {
	urlAddresses := make([]string, 0, len(urls))
	for _, url := range urls {
		urlAddresses = append(urlAddresses, url.Address)
	}
	return urlAddresses
}

type URL struct {
	ID              int32
	Address         string
	CategoryCode    CategoryCode
	BookmarkCount   int32
	NamedUserCount  int32
	PrivateUserRate float64
}
