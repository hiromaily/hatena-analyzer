package registry

import (
	"github.com/hiromaily/hatena-fake-detector/pkg/app"
	"github.com/hiromaily/hatena-fake-detector/pkg/logger"
)

type Registry interface {
	InitializeApp() (app.Application, error)
	Logger() logger.Logger
}
