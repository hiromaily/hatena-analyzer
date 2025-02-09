package registry

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/hiromaily/hatena-fake-detector/pkg/app"
	"github.com/hiromaily/hatena-fake-detector/pkg/envs"
	"github.com/hiromaily/hatena-fake-detector/pkg/handler"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
	"github.com/hiromaily/hatena-fake-detector/pkg/usecase"
)

type registry struct {
	envConf       *envs.Config
	appCode       app.AppCode
	commitID      string
	isCLI         bool
	targetHandler handler.Handler

	// common instance
	logger logger.Logger
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
		isCLI:    true, // CLI mode
	}
	reg.targetFunc()
	return &reg
}

func (r *registry) InitializeApp() (app.Application, error) {
	// CLI Application
	app := app.NewCLIApp(r.targetHandler)
	return app, nil
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
		r.envConf.InfluxdbURL,
		r.envConf.InfluxdbToken,
		r.envConf.InfluxdbBucket,
		r.envConf.InfluxdbOrg,
	)
}

///
/// Common instances
///

func (r *registry) newLogger() logger.Logger {
	if r.logger == nil {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		r.logger = logger.NewSlogLogger(
			slog.LevelInfo,
			hostname,
			r.appCode.String(),
			r.commitID,
		)
		//r.logger.Info("Logger initialized", "logLevel", logLevel.String())
	}
	return r.logger
}
