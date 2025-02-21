package entities

import (
	"errors"
	"strings"
)

type CategoryCode string

const (
	Unknown       CategoryCode = "unknown"
	All           CategoryCode = "all"
	General       CategoryCode = "general"
	Social        CategoryCode = "social"
	Economics     CategoryCode = "economics"
	Life          CategoryCode = "life"
	Knowledge     CategoryCode = "knowledge"
	IT            CategoryCode = "it"
	Fun           CategoryCode = "fun"
	Entertainment CategoryCode = "entertainment"
	Game          CategoryCode = "game"
)

func (c CategoryCode) String() string {
	return string(c)
}

func GetCategoryCodeList() []CategoryCode {
	return []CategoryCode{
		All, General, Social, Economics, Life, Knowledge, IT, Fun, Entertainment, Game,
	}
}

// convert to CategoryCode
func ToCategoryCode(s string) (CategoryCode, error) {
	switch CategoryCode(s) {
	case All, General, Social, Economics, Life, Knowledge, IT, Fun, Entertainment, Game:
		return CategoryCode(s), nil
	default:
		return "", errors.New("invalid category code")
	}
}

// extract category code from URL
func ExtractCategoryFromURL(urlStr string) (CategoryCode, error) {
	parts := strings.Split(urlStr, "/")

	if len(parts) == 0 {
		return "", errors.New("invalid URL format")
	}
	categoryPart := parts[len(parts)-1]

	return ToCategoryCode(categoryPart)
}
