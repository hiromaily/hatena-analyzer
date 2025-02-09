package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
	"github.com/joho/godotenv"
)

type Bookmark struct {
	Title     string `json:"title"`
	Count     int    `json:"count"`
	Users     map[string]User
	Timestamp time.Time
}

type User struct {
	Name        string `json:"name"`
	IsCommented bool   `json:"is_commented"`
	IsDeleted   bool   `json:"is_deleted"`
}

func main() {
	// .envファイルから環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 環境変数から値を取得
	influxdbURL := os.Getenv("INFLUXDB_URL")
	influxdbToken := os.Getenv("INFLUXDB_TOKEN")
	bucket := os.Getenv("INFLUXDB_BUCKET")
	org := os.Getenv("INFLUXDB_ORG")

	slogLogger := logger.NewSlogLogger(slog.LevelDebug, "localhost")

	// 初期化
	fetchUsecaser := usecase.NewFetchUsecase(slogLogger, influxdbURL, influxdbToken, bucket, org)
	err = fetchUsecaser.Execute(context.Background())
	if err != nil {
		slogLogger.Error("Failed to fetch bookmark data", "error", err)
	}
}
