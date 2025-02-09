package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Bookmark struct {
	Title string `json:"title"`
	Count int    `json:"count"`
	Users map[string]User
}

type User struct {
	Name        string `json:"name"`
	IsCommented bool   `json:"is_commented"`
}

func fetchBookmarkData(url string) (Bookmark, error) {
	apiUrl := "https://b.hatena.ne.jp/entry/json/" + url

	resp, err := http.Get(apiUrl)
	if err != nil {
		return Bookmark{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Title     string `json:"title"`
		Count     int    `json:"count"`
		Bookmarks []struct {
			User    string `json:"user"`
			Comment string `json:"comment"`
		} `json:"bookmarks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Bookmark{}, err
	}

	// ユーザー情報をマップに変換
	users := make(map[string]User)
	for _, bookmark := range data.Bookmarks {
		users[bookmark.User] = User{
			Name:        bookmark.User,
			IsCommented: bookmark.Comment != "",
		}
	}

	return Bookmark{
		Title: data.Title,
		Count: data.Count,
		Users: users,
	}, nil
}

func main() {
	urls := []string{
		"https://note.com/simplearchitect/n/nadc0bcdd5b3d",
		"https://note.com/simplearchitect/n/n871f29ffbfac",
		"https://note.com/simplearchitect/n/n86a95bc19b4c",
		"https://note.com/simplearchitect/n/nfd147540e3db",
	}

	for _, url := range urls {
		bookmark, err := fetchBookmarkData(url)
		if err != nil {
			fmt.Printf("Error fetching data for %s: %v\n", url, err)
			continue
		}

		fmt.Println("===================================================================")
		fmt.Printf("Title: %s\n", bookmark.Title)
		fmt.Printf("Count: %d\n", bookmark.Count)
		fmt.Printf("UserCount: %d\n", len(bookmark.Users))
		fmt.Printf("Users:\n")
		// for _, user := range bookmark.Users {
		// 	fmt.Printf(" - %s\n", user.Name)
		// }
		fmt.Println()
	}
}
