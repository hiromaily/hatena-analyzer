# SQL Linter

Linter by [sqlfluff](https://github.com/sqlfluff/sqlfluff) developed by Python

## Requirements

- Docker

## Install

```sh
make build-sqlfluff-image
```

## How to use

In case of working on `sql-linter` directory, prepare sql file in this directory.

```sh
# formatter
docker run --rm -v $(pwd):/workspace sqlfluff-local:latest fix /workspace/example.sql

# linter
docker run --rm -v $(pwd):/workspace sqlfluff-local:latest lint /workspace/example.sql
```

## References

- [Official Docs](https://docs.sqlfluff.com/en/stable/index.html)
