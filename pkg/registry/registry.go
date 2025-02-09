package registry

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/hiromaily/hatena-fake-detector/pkg/app"
	"github.com/hiromaily/hatena-fake-detector/pkg/envs"
	"github.com/hiromaily/hatena-fake-detector/pkg/fetcher"
	"github.com/hiromaily/hatena-fake-detector/pkg/handler"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/repository"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type registry struct {
	envConf       *envs.Config
	appCode       app.AppCode
	commitID      string
	isCLI         bool
	targetHandler handler.Handler

	// common instance
	logger          logger.Logger
	influxdbClient  influxdb2.Client
	bookmarkRepo    repository.BookmarkRepositorier
	bookmarkFetcher fetcher.BookmarkFetcher
}

func NewRegistry(
	envConf *envs.Config,
	appCode app.AppCode,
	commitID string,
) Registry {
	reg := registry{
		envConf:  envConf,
		appCode:  appCode,
		commitID: commitID,
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
	return nil, errors.New("Web Application is not implemented yet")
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
	}
	if r.targetHandler != nil {
		return
	}
	panic(fmt.Errorf("appCode is not found"))
}

///
/// handlers
///

// Health Handler
func (r *registry) newFetchHandler() handler.Handler {
	return handler.NewFetchCLIHandler(r.newLogger(), r.newFetchUsecase())
}

///
/// usecases
///

// must be called only once
func (r *registry) newFetchUsecase() usecase.FetchUsecaser {
	return usecase.NewFetchUsecase(
		r.newLogger(),
		r.newBookmarkRepository(),
		r.newBookmarkFetcher(),
	)
}

///
/// Common instances
///

func (r *registry) newLogger() logger.Logger {
	if r.logger == nil {
		r.logger = logger.NewSlogLogger(
			slog.LevelInfo, // default log level
			r.appCode.String(),
			r.commitID,
		)
		//r.logger.Info("Logger initialized", "logLevel", logLevel.String())
	}
	return r.logger
}

func (r *registry) newBookmarkRepository() repository.BookmarkRepositorier {
	if r.bookmarkRepo == nil {
		r.bookmarkRepo = repository.NewInfluxDBBookmarkRepository(
			r.newLogger(),
			r.newInfluxdbClient(),
			r.envConf.InfluxdbOrg,
			r.envConf.InfluxdbBucket,
		)
	}
	return r.bookmarkRepo
}

func (r *registry) newInfluxdbClient() influxdb2.Client {
	if r.influxdbClient == nil {
		r.influxdbClient = influxdb2.NewClient(r.envConf.InfluxdbURL, r.envConf.InfluxdbToken)
	}
	return r.influxdbClient
}

func (r *registry) newBookmarkFetcher() fetcher.BookmarkFetcher {
	if r.bookmarkFetcher == nil {
		r.bookmarkFetcher = fetcher.NewBookmarkFetcher(r.newLogger())
	}
	return r.bookmarkFetcher
}
