package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

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
	ginEngine *gin.Engine
	port      uint
}

func NewWebApp(ginEngine *gin.Engine, port uint) Application {
	// create web application
	if port == 0 {
		port = 8080
	}
	return &webApp{ginEngine: ginEngine, port: port}
}

func (c *webApp) Run() error {
	return c.ginEngine.Run(fmt.Sprintf(":%d", c.port))
}
