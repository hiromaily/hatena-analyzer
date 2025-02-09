package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hiromaily/hatena-fake-detector/pkg/entities"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type FetchUsecaser interface {
	Execute(ctx context.Context) error
}

type fetchUsecase struct {
	logger   logger.Logger
	dbClient influxdb2.Client
	bucket   string
	org      string
	urls     []string
	//repo              rdbrepo.RDBRepositorier
}

func NewFetchUsecase(
	logger logger.Logger,
	influxdbURL, influxdbToken, bucket, org string,
) *fetchUsecase {
	// InfluxDBクライアントを作成
	client := influxdb2.NewClient(influxdbURL, influxdbToken)

	// target URLリスト
	urls := []string{
		"https://note.com/simplearchitect/n/nadc0bcdd5b3d",
		"https://note.com/simplearchitect/n/n871f29ffbfac",
		"https://note.com/simplearchitect/n/n86a95bc19b4c",
		"https://note.com/simplearchitect/n/nfd147540e3db",
	}

	return &fetchUsecase{
		logger:   logger,
		dbClient: client,
		bucket:   bucket,
		org:      org,
		urls:     urls,
	}
}

func fetchBookmarkData(url string) (entities.Bookmark, error) {
	apiUrl := "https://b.hatena.ne.jp/entry/json/" + url

	resp, err := http.Get(apiUrl)
	if err != nil {
		return entities.Bookmark{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Title     string `json:"title"`
		Count     int    `json:"count"`
		Bookmarks []struct {
			User      string `json:"user"`
			Comment   string `json:"comment"`
			Timestamp string `json:"timestamp"`
		} `json:"bookmarks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return entities.Bookmark{}, err
	}

	users := make(map[string]entities.User)
	for _, bookmark := range data.Bookmarks {
		users[bookmark.User] = entities.User{
			Name:        bookmark.User,
			IsDeleted:   false,
			IsCommented: bookmark.Comment != "",
		}
	}

	return entities.Bookmark{
		Title:     data.Title,
		Count:     data.Count,
		Users:     users,
		Timestamp: time.Now(),
	}, nil
}

func writeBookmarkData(client influxdb2.Client, bucket, org, url string, bookmark entities.Bookmark) error {
	writeAPI := client.WriteAPIBlocking(org, bucket)
	tags := map[string]string{"url": url}
	fields := map[string]interface{}{
		"title": bookmark.Title,
		"count": bookmark.Count,
	}

	// Bookmarkデータポイントの作成
	point := write.NewPoint("bookmark", tags, fields, bookmark.Timestamp) // with timestamp

	err := writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		return err
	}

	// Userデータポイントの作成
	for _, user := range bookmark.Users {
		userFields := map[string]interface{}{
			"name":         user.Name,
			"is_commented": user.IsCommented,
			"is_deleted":   user.IsDeleted,
		}
		userPoint := write.NewPoint("user", tags, userFields, bookmark.Timestamp)
		err = writeAPI.WritePoint(context.Background(), userPoint)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadExistingData(client influxdb2.Client, bucket, org string, url string) (entities.Bookmark, error) {
	var bookmark entities.Bookmark
	bookmark.Users = make(map[string]entities.User)

	queryAPI := client.QueryAPI(org)
	query := fmt.Sprintf(`
	from(bucket: "%s")
	  |> range(start: 0)
	  |> filter(fn: (r) => r._measurement == "user" and r.url == "%s")
	  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	`, bucket, url)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return bookmark, err
	}

	for result.Next() {
		record := result.Record()
		// Debug: 詳細なログ出力を追加
		fmt.Printf("Record: %+v\n", record)

		// 必要なキーが存在するかをチェックし、適切な型にアサートする
		userName, ok := record.ValueByKey("name").(string)
		if !ok {
			fmt.Println("name field missing or not a string")
			continue
		}

		isDeleted, _ := record.ValueByKey("is_deleted").(bool)
		isCommented, _ := record.ValueByKey("is_commented").(bool)

		bookmark.Users[userName] = entities.User{
			Name:        userName,
			IsDeleted:   isDeleted,
			IsCommented: isCommented,
		}
	}
	if result.Err() != nil {
		return bookmark, result.Err()
	}

	return bookmark, nil
}

func (f *fetchUsecase) Execute(ctx context.Context) error {

	defer f.dbClient.Close()

	for _, url := range f.urls {
		// 既存データの読み込み
		existingBookmark, err := loadExistingData(f.dbClient, f.bucket, f.org, url)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				// 初回実行時の初期化
				existingBookmark.Users = make(map[string]entities.User)
			} else {
				fmt.Printf("Error loading existing data for %s: %v\n", url, err)
				continue
			}
		}

		// 既存ユーザーをすべて`isDeleted = true`に設定
		for userName := range existingBookmark.Users {
			existingBookmark.Users[userName] = entities.User{
				Name:        userName,
				IsDeleted:   true,
				IsCommented: existingBookmark.Users[userName].IsCommented,
			}
		}

		// 新しいデータを取得
		newBookmark, err := fetchBookmarkData(url)
		if err != nil {
			fmt.Printf("Error fetching data for %s: %v\n", url, err)
			continue
		}

		// 取得したユーザーで`isDeleted = false`に設定
		for userName, user := range newBookmark.Users {
			existingBookmark.Users[userName] = entities.User{
				Name:        userName,
				IsDeleted:   false,
				IsCommented: user.IsCommented,
			}
		}

		existingBookmark.Title = newBookmark.Title
		existingBookmark.Count = newBookmark.Count
		existingBookmark.Timestamp = newBookmark.Timestamp

		// データを保存
		err = writeBookmarkData(f.dbClient, f.bucket, f.org, url, existingBookmark)
		if err != nil {
			fmt.Printf("Error writing data to InfluxDB for %s: %v\n", url, err)
			continue
		}
		fmt.Printf("Data saved for URL: %s\n", url)

		// 表示
		fmt.Println("===================================================================")
		fmt.Printf("Title: %s\n", existingBookmark.Title)
		fmt.Printf("Count: %d\n", existingBookmark.Count)
		fmt.Printf("UserCount: %d\n", len(existingBookmark.Users))
		fmt.Printf("Users:\n")
		// for _, user := range existingBookmark.Users {
		// 	fmt.Printf(" - %s\n", user.Name)
		// }
		fmt.Println()
	}

	return nil
}
