package entities

import (
	"time"
)

type Bookmark struct {
	Title     string `json:"title"`
	Count     int    `json:"count"`
	Users     map[string]User
	Timestamp time.Time
}

type BookmarkSummary struct {
	Title     string `json:"title"`
	Count     int    `json:"count"`
	UserCount int    `json:"user_count"`
	Timestamp time.Time
}

type User struct {
	Name        string `json:"name"`
	IsCommented bool   `json:"is_commented"`
	IsDeleted   bool   `json:"is_deleted"`
}
