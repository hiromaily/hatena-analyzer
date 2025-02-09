package handler

import (
	"context"
)

type Handler interface {
	Handler(ctx context.Context) error
}
