package adapter

import (
	"github.com/hiromaily/hatena-analyzer/pkg/entities"
	"github.com/hiromaily/hatena-analyzer/pkg/storage/rdb/sqlcgen"
)

func DBUsersToEntityModel(users []sqlcgen.GetUsersByURLRow) []entities.RDBUser {
	var userModels []entities.RDBUser
	for _, user := range users {
		userModels = append(userModels, entities.RDBUser{
			UserName:      user.UserName,
			BookmarkCount: int(user.BookmarkCount.Int32),
		})
	}
	return userModels
}
