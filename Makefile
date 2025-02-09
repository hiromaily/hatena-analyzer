CURRENT_DIR := $(shell pwd)
# GO_VERSION=1.23.0
LINT_BIN=go run github.com/golangci/golangci-lint/cmd/golangci-lint
GOVULNCHECK_BIN=go run golang.org/x/vuln/cmd/govulncheck
GOMAJOR_BIN=go run github.com/icholy/gomajor
GOLINE_BIN=go run github.com/segmentio/golines
GOTESTSUM_BIN=go run gotest.tools/gotestsum

#------------------------------------------------------------------------------
# Tools for maintenance
#  毎回実行する必要はないが、定期的な確認に利用
#------------------------------------------------------------------------------
.PHONY: vulncheck
vulncheck:
	$(GOVULNCHECK_BIN) ./...

.PHONY: versioncheck
versioncheck:
	$(GOMAJOR_BIN) list

#------------------------------------------------------------------------------
# Lint
#------------------------------------------------------------------------------

# golangciの中で実行される`lll`の設定より短く指定する
.PHONY: linecheck
linecheck:
	$(GOLINE_BIN) -m 110 -w ./

.PHONY: lint
lint:
	$(LINT_BIN) run

.PHONY: lint-fix
lint-fix: linecheck
	$(LINT_BIN) run --fix

#------------------------------------------------------------------------------
# Execution
#------------------------------------------------------------------------------

.PHONY: run
run:
	go run ./cmd/fake-detector/
