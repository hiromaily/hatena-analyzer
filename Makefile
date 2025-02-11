CURRENT_DIR := $(shell pwd)
# GO_VERSION=1.23.0
LINT_BIN=go run github.com/golangci/golangci-lint/cmd/golangci-lint
GOVULNCHECK_BIN=go run golang.org/x/vuln/cmd/govulncheck
GOMAJOR_BIN=go run github.com/icholy/gomajor
GOLINE_BIN=go run github.com/segmentio/golines
GOTESTSUM_BIN=go run gotest.tools/gotestsum
SQLC_BIN=go run github.com/sqlc-dev/sqlc/cmd/sqlc

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
# Code generation
#------------------------------------------------------------------------------
.PHONY: copy-query
copy-query:
	cp ./docker/postgres/schema.sql ./sqlc/schemas/

.PHONY: gen-db-code
gen-db-code: clean-sqlc-gen-code
	$(SQLC_BIN) -f sqlc/sqlc.yml vet
	$(SQLC_BIN) -f sqlc/sqlc.yml generate

.PHONY: clean-gen-code
clean-sqlc-gen-code:
	rm -rf ./pkg/storage/rdb/sqlcgen

#------------------------------------------------------------------------------
# Execution
#------------------------------------------------------------------------------

.PHONY: run-fetch
run-fetch:
	go run ./cmd/fake-detector/ fetch

.PHONY: view-summary
view-summary:
	go run ./cmd/fake-detector/ view-summary

