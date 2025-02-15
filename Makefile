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
	cp ./docker/postgres/schema.sql ./tools/sqlc/schemas/
	cp ./docker/postgres/stored.sql ./tools/sqlc/schemas/

.PHONY: gen-db-code
gen-db-code: clean-sqlc-gen-code
	$(SQLC_BIN) -f tools/sqlc/sqlc.yml vet
	$(SQLC_BIN) -f tools/sqlc/sqlc.yml generate

.PHONY: clean-gen-code
clean-sqlc-gen-code:
	rm -rf ./pkg/storage/rdb/sqlcgen

.PHONY: gen-db-all
gen-db-all: copy-query gen-db-code

#------------------------------------------------------------------------------
# Execution
#------------------------------------------------------------------------------

# Fetch bookmarked entity from URLs and save data to the database
# saved data: bookmark URLs, bookmarked users and relations between them
.PHONY: fetch-bookmark
fetch-bookmark:
	go run ./cmd/fake-detector/ fetch-bookmark

# View time series data of the summary of bookmarked entity
.PHONY: view-summary
view-summary:
	go run ./cmd/fake-detector/ view-summary

# Update user's bookmark count
.PHONY: update-user-info
update-user-info:
	go run ./cmd/fake-detector/ update-user-info

# Run all executions
.PHONY: run-all
run-all: fetch-bookmark update-user-info
