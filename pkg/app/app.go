package app

import (
	"context"

	"github.com/hiromaily/hatena-analyzer/pkg/handler"
)

///
/// CLI Application
///

type cliApp struct {
	targetHandler handler.Handler
}

func NewCLIApp(handler handler.Handler) Application {
	return &cliApp{
		targetHandler: handler,
	}
}

func (c *cliApp) Run() error {
	return c.targetHandler.Handler(context.Background())
}

///
/// Web Application
///

type webApp struct {
}

func NewWebApp() Application {
	// create web application

	return &webApp{}
}

func (c *webApp) Run() error {
	return nil
}
