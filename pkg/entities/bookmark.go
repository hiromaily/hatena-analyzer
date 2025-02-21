package entities

import (
	"time"
)

type BookmarkUser struct {
	Name        string `json:"name"`
	IsCommented bool   `json:"is_commented"`
	IsDeleted   bool   `json:"is_deleted"`
}

type Bookmark struct {
	Title     string `json:"title"`
	Count     int    `json:"count"`
	Users     map[string]BookmarkUser
	Timestamp time.Time
}

func (b *Bookmark) CountDeletedUser() int {
	var count int
	for _, user := range b.Users {
		if user.IsDeleted {
			count++
		}
	}
	return count
}

type BookmarkSummary struct {
	Title            string `json:"title"`
	Count            int    `json:"count"`
	UserCount        int    `json:"user_count"`
	DeletedUserCount int    `json:"deleted_user_count"`
	Timestamp        time.Time
}
