# SQL Linter

Linter by [sqlfluff](https://github.com/sqlfluff/sqlfluff) developed by Python

## Requirements

- Docker

## Install

```sh
make build-sqlfluff-image
```

## How to use

`sql-linter`ディレクトリ内で作業する場合、ディレクトリ内に適当なsqlファイルを配置する。ここでは`sample.sql`とする。

```sh

# formatter
docker run --rm -v $(pwd):/workspace sqlfluff-local:latest fix /workspace/sample.sql

# linter
docker run --rm -v $(pwd):/workspace sqlfluff-local:latest lint /workspace/sample.sql
```

参考: [Makefile](./Makefile)、[batch/Makefile](../batch/Makefile): lint-queryターゲット

## References

- [Officialドキュメント](https://docs.sqlfluff.com/en/stable/index.html)
