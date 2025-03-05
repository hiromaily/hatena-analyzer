package entities

type URLIDAddress struct {
	ID      int32
	Address string
}

func FilterURLAddress(urls []URL) []string {
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
	Title           string
	BookmarkCount   int32
	NamedUserCount  int32
	PrivateUserRate float64
}

type AveragePrivateUserRate struct {
	CategoryCode           CategoryCode
	AveragePrivateUserRate float64
}

type LinkInfo struct {
	Href     string
	Category CategoryCode
	IsAll    bool
}

type LinkInfos []LinkInfo

func (l LinkInfos) Extract() ([]string, []CategoryCode, []bool) {
	urls := make([]string, len(l))
	categories := make([]CategoryCode, len(l))
	isAlls := make([]bool, len(l))
	for i, v := range l {
		urls[i] = v.Href
		categories[i] = v.Category
		isAlls[i] = v.IsAll
	}
	return urls, categories, isAlls
}
