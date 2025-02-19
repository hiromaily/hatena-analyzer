package registry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/hiromaily/hatena-fake-detector/pkg/app"
	"github.com/hiromaily/hatena-fake-detector/pkg/envs"
	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/handler"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/influxdb"
	"github.com/hiromaily/hatena-fake-detector/pkg/storage/rdb"
	"github.com/hiromaily/hatena-fake-detector/pkg/tracer"
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
	userRepo     repository.UserRepositorier

	// common instance
	logger              logger.Logger
	tracer              tracer.Tracer
	postgresClient      *rdb.SqlcPostgresClient
	influxdbClient      influxdb2.Client
	mongodbClient       *mongo.Client
	bookmarkFetcher     fetcher.BookmarkFetcher
	userBookmarkFetcher fetcher.UserBookmarkFetcher
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
	case r.appCode == app.AppCodeFetchBookmark:
		r.targetHandler = r.newFetchBookmarkHandler()
	case r.appCode == app.AppCodeViewSummary:
		r.targetHandler = r.newViewSummaryHanlder()
	case r.appCode == app.AppCodeUpdateUserInfo:
		r.targetHandler = r.newUpdateUserInfoHandler()
	}
	if r.targetHandler != nil {
		return
	}
	panic(fmt.Errorf("appCode is not found"))
}

///
/// handlers
///

func (r *registry) newFetchBookmarkHandler() handler.Handler {
	return handler.NewFetchBookmarkCLIHandler(r.newLogger(), r.newFetchBookmarkUsecase())
}

func (r *registry) newViewSummaryHanlder() handler.Handler {
	return handler.NewViewSummaryCLIHandler(r.newLogger(), r.newViewSummaryUsecase())
}

func (r *registry) newUpdateUserInfoHandler() handler.Handler {
	return handler.NewUpdateUserInfoCLIHandler(r.newLogger(), r.newUpdateUserInfoUsecase())
}

///
/// usecases
///

// must be called only once
func (r *registry) newFetchBookmarkUsecase() usecase.FetchBookmarkUsecaser {
	usecase, err := usecase.NewFetchBookmarkUsecase(
		r.newLogger(),
		r.newTracer(r.appCode.String()),
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
		r.newTracer(r.appCode.String()),
		r.newSummaryRepository(),
		r.urls,
	)
	if err != nil {
		panic(err)
	}
	return usecase
}

func (r *registry) newUpdateUserInfoUsecase() usecase.UpdateUserInfoUsecaser {
	usecase, err := usecase.NewUpdateUserInfoUsecase(
		r.newLogger(),
		r.newTracer(r.appCode.String()),
		r.newUserRepository(),
		r.newUserBookmarkFetcher(),
		r.envConf.MaxWorkers, // maxWorker
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
		// PosgreSQL implementation
		postgresBookmarkRepo := repository.NewRDBBookmarkRepository(
			r.newLogger(),
			r.newPostgresClient(),
		)

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
			postgresBookmarkRepo,
			influxdbBookmarkRepo,
			mongodbBookmarkRepo,
		)
	}
	return r.bookmarkRepo
}

func (r *registry) newSummaryRepository() repository.SummaryRepositorier {
	if r.summaryRepo == nil {
		// InfluxDB implementation
		influxDBSummaryRepo := repository.NewInfluxDBSummaryRepository(
			r.newLogger(),
			r.newInfluxdbClient(),
			r.envConf.InfluxdbOrg,
			r.envConf.InfluxdbBucket,
		)
		// PosgreSQL implementation
		rdbSummaryRepo := repository.NewRDBSummaryRepository(
			r.newLogger(),
			r.newPostgresClient(),
		)
		r.summaryRepo = repository.NewSummaryRepository(
			r.newLogger(),
			rdbSummaryRepo,
			influxDBSummaryRepo,
		)
	}
	return r.summaryRepo
}

func (r *registry) newUserRepository() repository.UserRepositorier {
	if r.userRepo == nil {
		// PosgreSQL implementation
		r.userRepo = repository.NewRDBUserRepository(
			r.newLogger(),
			r.newPostgresClient(),
		)
	}
	return r.userRepo
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

func (r *registry) newTracer(tracerName string) tracer.Tracer {
	if r.tracer == nil {
		var err error
		serviceName := r.envConf.TracerServiceName
		version := r.envConf.TracerVersion
		sampler := sdktrace.AlwaysSample()
		tracerMode := tracer.ValidateTracerEnv(r.envConf.Tracer)

		switch tracerMode {
		case tracer.TracerModeNOOP:
			r.tracer = tracer.NewNoopProvider()
		case tracer.TracerModeJaegerHTTP:
			host := "localhost:4318"
			r.tracer, err = tracer.NewJaegerHTTPProvider(host, serviceName, tracerName, version, sampler)
		case tracer.TracerModeJaegerGRPC:
			host := "localhost:4317"
			r.tracer, err = tracer.NewJaegerGRPCProvider(host, serviceName, tracerName, version, sampler)
		// case tracer.TracerModeDataDog:
		// 	// datadog
		// 	r.tracer, err = tracer.NewDatadogOtelProvider(tracerName, version, isDebug)
		// 	r.tracer = tracer.NewDatadogTracer()
		default:
			err = errors.New("environment variable: Tracer is invalid")
		}
		if err != nil {
			panic(err)
		}
	}
	return r.tracer
}

func (r *registry) newPostgresClient() *rdb.SqlcPostgresClient {
	if r.postgresClient == nil {
		pgClient, err := rdb.NewSqlcPostgresClient(
			context.Background(),
			r.envConf.PostgresURL,
			r.envConf.DBMaxConnections,
		)
		if err != nil {
			panic(err)
		}
		r.postgresClient = pgClient
	}
	return r.postgresClient
}

func (r *registry) newInfluxdbClient() influxdb2.Client {
	if r.influxdbClient == nil {
		r.influxdbClient = influxdb2.NewClient(r.envConf.InfluxdbURL, r.envConf.InfluxdbToken)
		// ping
		err := influxdb.Ping(r.influxdbClient, r.envConf.InfluxdbOrg, r.envConf.InfluxdbBucket)
		if err != nil {
			panic(err)
		}
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

func (r *registry) newUserBookmarkFetcher() fetcher.UserBookmarkFetcher {
	if r.userBookmarkFetcher == nil {
		r.userBookmarkFetcher = fetcher.NewBookmarkUserFetcher(r.newLogger())
	}
	return r.userBookmarkFetcher
}
