// tools.go
//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/icholy/gomajor"
	_ "github.com/segmentio/golines"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "golang.org/x/vuln/cmd/govulncheck"
	_ "gotest.tools/gotestsum"
)
