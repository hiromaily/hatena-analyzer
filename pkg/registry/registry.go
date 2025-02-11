package registry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/hiromaily/hatena-fake-detector/pkg/app"
	"github.com/hiromaily/hatena-fake-detector/pkg/envs"
	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/handler"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

type registry struct {
	envConf  *envs.Config
	appCode  app.AppCode
	commitID string
	urls     []string

	isCLI         bool
	targetHandler handler.Handler

	// repositories
	bookmarkRepo repository.BookmarkRepositorier
	summaryRepo  repository.SummaryRepositorier

	// common instance
	logger          logger.Logger
	postgresClient  *pgx.Conn
	influxdbClient  influxdb2.Client
	mongodbClient   *mongo.Client
	bookmarkFetcher fetcher.BookmarkFetcher
}

func NewRegistry(
	envConf *envs.Config,
	appCode app.AppCode,
	commitID string,
	urls []string,
) Registry {
	reg := registry{
		envConf:  envConf,
		appCode:  appCode,
		commitID: commitID,
		urls:     urls,
		isCLI:    appCode != app.AppCodeWeb, // CLI mode
	}
	reg.targetFunc()
	return &reg
}

func (r *registry) InitializeApp() (app.Application, error) {
	if r.isCLI {
		// CLI Application
		app := app.NewCLIApp(r.targetHandler)
		return app, nil
	}
	return nil, errors.New("web application is not implemented yet")
}

func (r *registry) Logger() logger.Logger {
	return r.newLogger()
}

func (r *registry) targetFunc() {
	if !r.isCLI {
		return
	}

	switch {
	case r.appCode == app.AppCodeFetch:
		r.targetHandler = r.newFetchHandler()
	case r.appCode == app.AppCodeViewSummary:
		r.targetHandler = r.newViewSummaryHanlder()
	}
	if r.targetHandler != nil {
		return
	}
	panic(fmt.Errorf("appCode is not found"))
}

///
/// handlers
///

func (r *registry) newFetchHandler() handler.Handler {
	return handler.NewFetchCLIHandler(r.newLogger(), r.newFetchUsecase())
}

func (r *registry) newViewSummaryHanlder() handler.Handler {
	return handler.NewViewSummaryCLIHandler(r.newLogger(), r.newViewSummaryUsecase())
}

///
/// usecases
///

// must be called only once
func (r *registry) newFetchUsecase() usecase.FetchUsecaser {
	usecase, err := usecase.NewFetchUsecase(
		r.newLogger(),
		r.newBookmarkRepository(),
		r.newBookmarkFetcher(),
		r.urls,
	)
	if err != nil {
		panic(err)
	}
	return usecase
}

func (r *registry) newViewSummaryUsecase() usecase.ViewSummaryUsecaser {
	usecase, err := usecase.NewViewSummaryUsecase(
		r.newLogger(),
		r.newSummaryRepository(),
		r.urls,
	)
	if err != nil {
		panic(err)
	}
	return usecase
}

///
/// Repositories
///

func (r *registry) newBookmarkRepository() repository.BookmarkRepositorier {
	if r.bookmarkRepo == nil {
		// TODO: Posgres implementation

		// InfluxDB implementation
		influxdbBookmarkRepo := repository.NewInfluxDBBookmarkRepository(
			r.newLogger(),
			r.newInfluxdbClient(),
			r.envConf.InfluxdbOrg,
			r.envConf.InfluxdbBucket,
		)
		// MongoDB implementation
		mongodbBookmarkRepo := repository.NewMongoDBBookmarkRepository(
			r.newLogger(),
			r.newMongodbClient(),
			r.envConf.MongodbDB,
			r.envConf.MongodbCollection,
		)

		r.bookmarkRepo = repository.NewBookmarkRepository(
			r.newLogger(),
			influxdbBookmarkRepo,
			mongodbBookmarkRepo,
		)
	}
	return r.bookmarkRepo
}

func (r *registry) newSummaryRepository() repository.SummaryRepositorier {
	if r.summaryRepo == nil {
		// InfluxDB implementation
		r.summaryRepo = repository.NewInfluxDBSummaryRepository(
			r.newLogger(),
			r.newInfluxdbClient(),
			r.envConf.InfluxdbOrg,
			r.envConf.InfluxdbBucket,
		)
	}
	return r.summaryRepo
}

///
/// Common instances
///

func (r *registry) newLogger() logger.Logger {
	if r.logger == nil {
		logLevel := slog.LevelInfo // default log level
		if r.envConf.IsDebug {
			logLevel = slog.LevelDebug
		}
		r.logger = logger.NewSlogLogger(
			logLevel,
			r.appCode.String(),
			r.commitID,
		)
		// r.logger.Info("Logger initialized", "logLevel", logLevel.String())
	}
	return r.logger
}

func (r *registry) newPostgresClient() *pgx.Conn {
	if r.postgresClient == nil {
		conn, err := pgx.Connect(context.Background(), r.envConf.PostgresURL)
		if err != nil {
			panic(err)
		}
		r.postgresClient = conn
	}
	return r.postgresClient
}

func (r *registry) newInfluxdbClient() influxdb2.Client {
	if r.influxdbClient == nil {
		r.influxdbClient = influxdb2.NewClient(r.envConf.InfluxdbURL, r.envConf.InfluxdbToken)
	}
	return r.influxdbClient
}

func (r *registry) newMongodbClient() *mongo.Client {
	if r.mongodbClient == nil {
		clientOptions := options.Client().ApplyURI(r.envConf.MongodbURL)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			panic(err)
		}
		r.mongodbClient = client
	}
	return r.mongodbClient
}

func (r *registry) newBookmarkFetcher() fetcher.BookmarkFetcher {
	if r.bookmarkFetcher == nil {
		r.bookmarkFetcher = fetcher.NewBookmarkFetcher(r.newLogger())
	}
	return r.bookmarkFetcher
}
