package registry

import (
	"github.com/hiromaily/hatena-analyzer/pkg/app"
	"github.com/hiromaily/hatena-analyzer/pkg/logger"
)

type Registry interface {
	InitializeApp() (app.Application, error)
	Logger() logger.Logger
}
