package handler

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Handler(ctx context.Context) error
	WebHandler(c *gin.Context)
}
