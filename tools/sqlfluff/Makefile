SQL_LINT_IMAGE_NAME = sqlfluff-local:latest

PWD := $(shell pwd)

.PHONY: build-sqlfluff-image
build-sqlfluff-image:
	docker build --progress=plain -t $(SQL_LINT_IMAGE_NAME) .

# format 実行例
# e.g. make format SQL=postgres_example.sql
.PHONY: format-sql
format-sql:
	docker run --rm -v $(PWD):/workspace $(SQL_LINT_IMAGE_NAME) fix /workspace/$(SQL)

# lint 実行例
# e.g make lint SQL=postgres_example.sql
.PHONY: lint-sql
lint-sql:
	docker run --rm -v $(PWD):/workspace $(SQL_LINT_IMAGE_NAME) lint /workspace/$(SQL)
