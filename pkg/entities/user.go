package entities

type RDBUser struct {
	UserName      string
	BookmarkCount int
}

func PrivateUserRate(totalCount, userCount int) float64 {
	// private user rate
	return float64(totalCount-userCount) / float64(totalCount) * 100
}
