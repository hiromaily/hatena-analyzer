package registry

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/hiromaily/hatena-analyzer/pkg/app"
	"github.com/hiromaily/hatena-analyzer/pkg/args"
	"github.com/hiromaily/hatena-analyzer/pkg/entities"
	"github.com/hiromaily/hatena-analyzer/pkg/envs"
	"github.com/hiromaily/hatena-analyzer/pkg/fetcher"
	"github.com/hiromaily/hatena-analyzer/pkg/handler"
	"github.com/hiromaily/hatena-analyzer/pkg/logger"
	"github.com/hiromaily/hatena-analyzer/pkg/repository"
	"github.com/hiromaily/hatena-analyzer/pkg/storage/influxdb"
	"github.com/hiromaily/hatena-analyzer/pkg/storage/mongodb"
	"github.com/hiromaily/hatena-analyzer/pkg/storage/rdb"
	"github.com/hiromaily/hatena-analyzer/pkg/tracer"
	"github.com/hiromaily/hatena-analyzer/pkg/usecase"
)

type registry struct {
	envConf  *envs.Config
	appCode  app.AppCode
	commitID string
	args     *args.Args

	isCLI bool

	// repositories
	closerRepo          repository.CloserRepositorier
	fetchBookmarkRepo   repository.FetchBookmarkRepositorier
	fetchURLRepo        repository.FetchURLRepositorier
	fetchUserRepo       repository.FetchUserRepositorier
	timeSeriesRepo      repository.TimeSeriesRepositorier
	bookmarkDetailsRepo repository.BookmarkDetailsRepositorier
	summaryRepo         repository.SummaryRepositorier

	// db clients
	postgresClient  *rdb.SqlcPostgresClient
	influxdbClient  influxdb2.Client
	mongodbClient   *mongo.Client
	postgresQueries *rdb.PostgreQueries
	influxDBQueries *influxdb.InfluxDBQueries
	mongoDBQueries  *mongodb.MongoDBQueries
	// fetchers
	entityJSONFetcher        fetcher.EntityJSONFetcher
	userBookmarkCountFetcher fetcher.UserBookmarkCountFetcher
	pageURLFetcher           fetcher.HatenaPageURLFetcher
	// common instance
	logger logger.Logger
	tracer tracer.Tracer
}

func NewRegistry(
	envConf *envs.Config,
	appCode app.AppCode,
	commitID string,
	args *args.Args,
) Registry {
	reg := registry{
		envConf:  envConf,
		appCode:  appCode,
		commitID: commitID,
		args:     args,
		isCLI:    appCode != app.AppCodeWeb, // CLI mode
	}
	return &reg
}

func (r *registry) InitializeApp() (app.Application, error) {
	_, err := r.newCloserRepository()
	if err != nil {
		return nil, err
	}

	if r.isCLI {
		// CLI Application
		handler, err := r.createCLIHandler()
		if err != nil {
			return nil, err
		}
		app := app.NewCLIApp(handler)
		return app, nil
	}
	// Web Application
	ginEngine := gin.Default()
	if err := r.createWebHandler(ginEngine); err != nil {
		return nil, err
	}
	app := app.NewWebApp(ginEngine, r.args.WebCommand.Port)
	return app, nil
}

func (r *registry) Logger() logger.Logger {
	return r.newLogger()
}

func (r *registry) Close() error {
	if r.closerRepo != nil {
		r.closerRepo.Close(context.Background())
	}
	return nil
}

func (r *registry) createCLIHandler() (handler.Handler, error) {
	var handler handler.Handler
	var err error

	switch {
	case r.appCode == app.AppCodeFetchHatenaPageURLs:
		handler, err = r.newFetchHatenaPageURLsHandler()
	case r.appCode == app.AppCodeFetchBookmarkEntities:
		handler, err = r.newFetchBookmarkHandler()
	case r.appCode == app.AppCodeFetchUserBookmarkCount:
		handler, err = r.newFetchUserBookmarkCountHandler()
	case r.appCode == app.AppCodeViewTimeSeries:
		handler, err = r.newViewTimeSeriesHanlder()
	case r.appCode == app.AppCodeViewBookmarkDetails:
		handler, err = r.newViewBookmarkDetailsHanlder()
	case r.appCode == app.AppCodeViewSummary:
		handler, err = r.newViewSummaryHanlder()
	}
	if err != nil {
		return nil, err
	}
	if handler == nil {
		return nil, errors.New("appCode is not found")
	}
	return handler, nil
}

func (r *registry) createWebHandler(ginEngine *gin.Engine) error {
	v1Router := ginEngine.Group("/api/v1")

	handler, err := r.newFetchHatenaPageURLsHandler()
	if err != nil {
		return err
	}
	v1Router.GET("/fetch-page-url", handler.WebHandler)

	handler, err = r.newFetchBookmarkHandler()
	if err != nil {
		return err
	}
	v1Router.GET("/fetch-bookmark", handler.WebHandler)

	handler, err = r.newFetchUserBookmarkCountHandler()
	if err != nil {
		return err
	}
	v1Router.GET("/fetch-user-bookmark-count", handler.WebHandler)

	handler, err = r.newViewTimeSeriesHanlder()
	if err != nil {
		return err
	}
	v1Router.GET("/view-time-series", handler.WebHandler)

	handler, err = r.newViewBookmarkDetailsHanlder()
	if err != nil {
		return err
	}
	v1Router.GET("/view-bookmark-details", handler.WebHandler)

	handler, err = r.newViewSummaryHanlder()
	if err != nil {
		return err
	}
	v1Router.GET("/view-summary", handler.WebHandler)

	return nil
}

///
/// CLI handlers
///

func (r *registry) newFetchHatenaPageURLsHandler() (handler.Handler, error) {
	usecaser, err := r.newFetchHatenaPageURLsUsecase()
	if err != nil {
		return nil, err
	}
	if r.isCLI {
		return handler.NewFetchHatenaPageURLsCLIHandler(r.newLogger(), usecaser), nil
	}
	return handler.NewFetchHatenaPageURLsWebHandler(r.newLogger(), usecaser), nil
}

func (r *registry) newFetchBookmarkHandler() (handler.Handler, error) {
	usecaser, err := r.newFetchBookmarkUsecase()
	if err != nil {
		return nil, err
	}
	if r.isCLI {
		// retrieve args
		var urls []string
		if r.args.FetchBookmarkEntitiesCommand.URLs != "" {
			urls = strings.Split(r.args.FetchBookmarkEntitiesCommand.URLs, ",")
			r.newLogger().Info("given URLs", "urls", urls, "len", len(urls))
		}
		return handler.NewFetchBookmarkCLIHandler(
			r.newLogger(), usecaser,
			urls, r.args.FetchBookmarkEntitiesCommand.Verbose,
		), nil
	}
	return handler.NewFetchBookmarkWebHandler(
		r.newLogger(), usecaser,
	), nil
}

func (r *registry) newFetchUserBookmarkCountHandler() (handler.Handler, error) {
	usecaser, err := r.newFetchUserBookmarkCountUsecase()
	if err != nil {
		return nil, err
	}
	if r.isCLI {
		// retrieve args
		var urls []string
		if r.args.FetchUserBookmarkCountCommand.URLs != "" {
			urls = strings.Split(r.args.FetchUserBookmarkCountCommand.URLs, ",")
			r.newLogger().Info("given URLs", "urls", urls, "len", len(urls))
		}
		return handler.NewFetchUserBookmarkCountCLIHandler(r.newLogger(), usecaser, urls), nil
	}
	return handler.NewFetchUserBookmarkCountWebHandler(r.newLogger(), usecaser), nil
}

func (r *registry) newViewTimeSeriesHanlder() (handler.Handler, error) {
	usecaser, err := r.newViewTimeSeriesUsecase()
	if err != nil {
		return nil, err
	}
	if r.isCLI {
		// retrieve args
		var urls []string
		if r.args.ViewTimeSeriesCommand.URLs != "" {
			urls = strings.Split(r.args.ViewTimeSeriesCommand.URLs, ",")
			r.newLogger().Info("given URLs", "urls", urls, "len", len(urls))
		}
		return handler.NewViewTimeSeriesCLIHandler(r.newLogger(), usecaser, urls), nil
	}
	return handler.NewViewTimeSeriesWebHandler(r.newLogger(), usecaser), nil
}

func (r *registry) newViewBookmarkDetailsHanlder() (handler.Handler, error) {
	usecaser, err := r.newViewBookmarkDetailsUsecase()
	if err != nil {
		return nil, err
	}
	if r.isCLI {
		// retrieve args
		var urls []string
		if r.args.ViewBookmarkDetailsCommand.URLs != "" {
			urls = strings.Split(r.args.ViewBookmarkDetailsCommand.URLs, ",")
			r.newLogger().Info("given URLs", "urls", urls, "len", len(urls))
		}
		return handler.NewViewBookmarkDetailsCLIHandler(r.newLogger(), usecaser, urls), nil
	}
	return handler.NewViewBookmarkDetailsWebHandler(r.newLogger(), usecaser), nil
}

func (r *registry) newViewSummaryHanlder() (handler.Handler, error) {
	// retrieve args
	usecaser, err := r.newViewSummaryUsecase()
	if err != nil {
		return nil, err
	}
	if r.isCLI {
		var urls []string
		if r.args.ViewSummaryCommand.URLs != "" {
			urls = strings.Split(r.args.ViewSummaryCommand.URLs, ",")
			r.newLogger().Info("given URLs", "urls", urls, "len", len(urls))
		}
		return handler.NewViewSummaryCLIHandler(
			r.newLogger(),
			usecaser,
			urls,
			r.args.ViewSummaryCommand.Threshold,
		), nil
	}
	return handler.NewViewSummaryWebHandler(r.newLogger(), usecaser), nil
}

///
/// usecases
///

// must be called only once

func (r *registry) newFetchHatenaPageURLsUsecase() (usecase.FetchHatenaPageURLsUsecaser, error) {
	tracer, err := r.newTracer(r.appCode.String())
	if err != nil {
		return nil, err
	}
	urlRepo, err := r.newURLRepository()
	if err != nil {
		return nil, err
	}
	usecase, err := usecase.NewFetchHatenaPageURLsUsecase(
		r.newLogger(),
		tracer,
		urlRepo,
		r.newPageURLFetcher(),
		entities.Unknown, // TODO: pass from cli args
	)
	if err != nil {
		return nil, err
	}
	return usecase, nil
}

func (r *registry) newFetchBookmarkUsecase() (usecase.FetchBookmarkUsecaser, error) {
	tracer, err := r.newTracer(r.appCode.String())
	if err != nil {
		return nil, err
	}
	bookmarkRepo, err := r.newBookmarkRepository()
	if err != nil {
		return nil, err
	}
	usecase, err := usecase.NewFetchBookmarkUsecase(
		r.newLogger(),
		tracer,
		bookmarkRepo,
		r.newBookmarkFetcher(),
		r.envConf.MaxWorkers, // maxWorker
	)
	if err != nil {
		return nil, err
	}
	return usecase, nil
}

func (r *registry) newViewTimeSeriesUsecase() (usecase.ViewTimeSeriesUsecaser, error) {
	tracer, err := r.newTracer(r.appCode.String())
	if err != nil {
		return nil, err
	}
	timeSeriesRepo, err := r.newTimeSeriesRepository()
	if err != nil {
		return nil, err
	}
	usecase, err := usecase.NewViewTimeSeriesUsecase(
		r.newLogger(),
		tracer,
		timeSeriesRepo,
	)
	if err != nil {
		return nil, err
	}
	return usecase, nil
}

func (r *registry) newViewBookmarkDetailsUsecase() (usecase.ViewBookmarkDetailsUsecaser, error) {
	tracer, err := r.newTracer(r.appCode.String())
	if err != nil {
		return nil, err
	}
	bookmarkDetailsRepo, err := r.newBookmarkDetailsRepository()
	if err != nil {
		return nil, err
	}
	usecase, err := usecase.NewViewBookmarkDetailsUsecase(
		r.newLogger(),
		tracer,
		bookmarkDetailsRepo,
	)
	if err != nil {
		return nil, err
	}
	return usecase, nil
}

func (r *registry) newViewSummaryUsecase() (usecase.ViewSummaryUsecaser, error) {
	tracer, err := r.newTracer(r.appCode.String())
	if err != nil {
		return nil, err
	}
	summaryRepo, err := r.newSummaryRepository()
	if err != nil {
		return nil, err
	}
	usecase, err := usecase.NewViewSummaryUsecase(
		r.newLogger(),
		tracer,
		summaryRepo,
	)
	if err != nil {
		return nil, err
	}
	return usecase, nil
}

func (r *registry) newFetchUserBookmarkCountUsecase() (usecase.FetchUserBookmarkCountUsecaser, error) {
	tracer, err := r.newTracer(r.appCode.String())
	if err != nil {
		return nil, err
	}
	userRepo, err := r.newUserRepository()
	if err != nil {
		return nil, err
	}
	usecase, err := usecase.NewFetchUserBookmarkCountUsecase(
		r.newLogger(),
		tracer,
		userRepo,
		r.newUserBookmarkCountFetcher(),
		r.envConf.MaxWorkers, // maxWorker
	)
	if err != nil {
		return nil, err
	}
	return usecase, nil
}

///
/// Repositories
///

func (r *registry) newCloserRepository() (repository.CloserRepositorier, error) {
	pgQuery, err := r.newPostgresQueries()
	if err != nil {
		return nil, err
	}
	influxdbQuery, err := r.newInfluxDBQueries()
	if err != nil {
		return nil, err
	}
	mongodbQuery, err := r.newMongoDBQueries()
	if err != nil {
		return nil, err
	}
	if r.closerRepo == nil {
		r.closerRepo = repository.NewCloserRepository(
			r.newLogger(),
			pgQuery,
			influxdbQuery,
			mongodbQuery,
		)
	}
	return r.closerRepo, nil
}

func (r *registry) newBookmarkRepository() (repository.FetchBookmarkRepositorier, error) {
	pgQuery, err := r.newPostgresQueries()
	if err != nil {
		return nil, err
	}
	influxdbQuery, err := r.newInfluxDBQueries()
	if err != nil {
		return nil, err
	}
	mongodbQuery, err := r.newMongoDBQueries()
	if err != nil {
		return nil, err
	}
	if r.fetchBookmarkRepo == nil {
		r.fetchBookmarkRepo = repository.NewFetchBookmarkRepository(
			r.newLogger(),
			pgQuery,
			influxdbQuery,
			mongodbQuery,
		)
	}
	return r.fetchBookmarkRepo, nil
}

func (r *registry) newTimeSeriesRepository() (repository.TimeSeriesRepositorier, error) {
	influxdbQuery, err := r.newInfluxDBQueries()
	if err != nil {
		return nil, err
	}
	if r.timeSeriesRepo == nil {
		r.timeSeriesRepo = repository.NewTimeSeriesRepository(
			r.newLogger(),
			influxdbQuery,
		)
	}
	return r.timeSeriesRepo, nil
}

func (r *registry) newBookmarkDetailsRepository() (repository.BookmarkDetailsRepositorier, error) {
	pgQuery, err := r.newPostgresQueries()
	if err != nil {
		return nil, err
	}
	if r.bookmarkDetailsRepo == nil {
		r.bookmarkDetailsRepo = repository.NewBookmarkDetailsRepository(
			r.newLogger(),
			pgQuery,
		)
	}
	return r.bookmarkDetailsRepo, nil
}

func (r *registry) newSummaryRepository() (repository.SummaryRepositorier, error) {
	pgQuery, err := r.newPostgresQueries()
	if err != nil {
		return nil, err
	}
	if r.summaryRepo == nil {
		r.summaryRepo = repository.NewSummaryRepository(
			r.newLogger(),
			pgQuery,
		)
	}
	return r.summaryRepo, nil
}

func (r *registry) newUserRepository() (repository.FetchUserRepositorier, error) {
	pgQuery, err := r.newPostgresQueries()
	if err != nil {
		return nil, err
	}
	if r.fetchUserRepo == nil {
		r.fetchUserRepo = repository.NewFetchUserRepository(
			r.newLogger(),
			pgQuery,
		)
	}
	return r.fetchUserRepo, nil
}

func (r *registry) newURLRepository() (repository.FetchURLRepositorier, error) {
	pgQuery, err := r.newPostgresQueries()
	if err != nil {
		return nil, err
	}
	if r.fetchURLRepo == nil {
		r.fetchURLRepo = repository.NewFetchURLRepository(
			r.newLogger(),
			pgQuery,
		)
	}
	return r.fetchURLRepo, nil
}

///
/// DB Clients
///

func (r *registry) newPostgresQueries() (*rdb.PostgreQueries, error) {
	pgClient, err := r.newPostgresClient()
	if err != nil {
		return nil, err
	}
	if r.postgresQueries == nil {
		r.postgresQueries = rdb.NewPostgreQueries(
			r.newLogger(),
			pgClient,
		)
	}
	return r.postgresQueries, nil
}

func (r *registry) newInfluxDBQueries() (*influxdb.InfluxDBQueries, error) {
	influxdbClient, err := r.newInfluxdbClient()
	if err != nil {
		return nil, err
	}
	if r.influxDBQueries == nil {
		r.influxDBQueries = influxdb.NewInfluxDBQueries(
			r.newLogger(),
			influxdbClient,
			r.envConf.InfluxdbOrg,
			r.envConf.InfluxdbBucket,
		)
	}
	return r.influxDBQueries, nil
}

func (r *registry) newMongoDBQueries() (*mongodb.MongoDBQueries, error) {
	mongodbClient, err := r.newMongodbClient()
	if err != nil {
		return nil, err
	}
	if r.mongoDBQueries == nil {
		r.mongoDBQueries = mongodb.NewMongoDBQueries(
			r.newLogger(),
			mongodbClient,
			r.envConf.MongodbDB,
			r.envConf.MongodbCollection,
		)
	}
	return r.mongoDBQueries, nil
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
		switch r.envConf.Logger {
		case "console":
			// console logger
			r.logger = logger.NewSlogConsoleLogger(
				logLevel,
			)
		case "json":
			// JSON logger
			r.logger = logger.NewSlogJSONLogger(
				logLevel,
				r.appCode.String(),
				r.commitID,
			)
		default:
			// default
			r.logger = logger.NewNoopLogger()
		}
		// r.logger.Info("Logger initialized", "logLevel", logLevel.String())
	}
	return r.logger
}

func (r *registry) newTracer(tracerName string) (tracer.Tracer, error) {
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
			return nil, err
		}
	}
	return r.tracer, nil
}

func (r *registry) newPostgresClient() (*rdb.SqlcPostgresClient, error) {
	if r.postgresClient == nil {
		pgClient, err := rdb.NewSqlcPostgresClient(
			context.Background(),
			r.envConf.PostgresURL,
			r.envConf.DBMaxConnections,
		)
		if err != nil {
			return nil, err
		}
		r.postgresClient = pgClient
	}
	return r.postgresClient, nil
}

func (r *registry) newInfluxdbClient() (influxdb2.Client, error) {
	if r.influxdbClient == nil {
		r.influxdbClient = influxdb2.NewClient(r.envConf.InfluxdbURL, r.envConf.InfluxdbToken)
		// ping
		err := influxdb.Ping(r.influxdbClient, r.envConf.InfluxdbOrg, r.envConf.InfluxdbBucket)
		if err != nil {
			return nil, err
		}
	}
	return r.influxdbClient, nil
}

func (r *registry) newMongodbClient() (*mongo.Client, error) {
	if r.mongodbClient == nil {
		clientOptions := options.Client().ApplyURI(r.envConf.MongodbURL)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			return nil, err
		}
		r.mongodbClient = client
	}
	return r.mongodbClient, nil
}

func (r *registry) newBookmarkFetcher() fetcher.EntityJSONFetcher {
	if r.entityJSONFetcher == nil {
		r.entityJSONFetcher = fetcher.NewEntityJSONFetcher(r.newLogger())
	}
	return r.entityJSONFetcher
}

func (r *registry) newUserBookmarkCountFetcher() fetcher.UserBookmarkCountFetcher {
	if r.userBookmarkCountFetcher == nil {
		r.userBookmarkCountFetcher = fetcher.NewUserBookmarkCountFetcher(r.newLogger())
	}
	return r.userBookmarkCountFetcher
}

func (r *registry) newPageURLFetcher() fetcher.HatenaPageURLFetcher {
	if r.pageURLFetcher == nil {
		r.pageURLFetcher = fetcher.NewHatenaPageURLFetcher(r.newLogger())
	}
	return r.pageURLFetcher
}
