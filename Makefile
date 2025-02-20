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

# must be shorter than `lll` settings in golangci-lint
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

# Fetch page urls from the Hatena pages
.PHONY: fetch-page-urls
fetch-page-urls:
	go run ./cmd/fake-detector/ fetch-hatena-page-urls

# Fetch bookmark users, title, count from page of given URL and save data to DB
.PHONY: fetch-bookmark
fetch-bookmark:
	go run ./cmd/fake-detector/ fetch-bookmark

# Fetch user's bookmark count
.PHONY: fetch-user-bm-count
fetch-user-bm-count:
	go run ./cmd/fake-detector/ fetch-user-bm-count

# View time series data of the summary of bookmarked entity
.PHONY: view-summary
view-summary:
	go run ./cmd/fake-detector/ view-summary

# Run all executions
.PHONY: run-all
run-all: fetch-bookmark fetch-user-bm-count view-summary
