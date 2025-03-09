# hatena-analyzer

Analyze bookmarked entity on [Hatena](https://b.hatena.ne.jp/hotentry/all). It works both CLI mode and Web server mode.

## Requirements

- Golang 1.23+
- Docker

## Commands

```sh
# Help
go run ./cmd/analyzer/ -h

# Build
go build -v -o ${GOPATH}/bin/hatena-analyzer ./cmd/analyzer/
```

### use as CLI

- `fetch-hatena-page-urls`: Fetched listed urls from Hatena page
- `fetch-bookmark`: Fetch bookmark entity information from url and save data to the database
- `fetch-user-bm-count`: Fetch user's bookmark count
- `view-timeseries`: View time series of bookmarked entity
- `view-bookmark-details`: View details of bookmarked entity
- `view-summary`: View summary of bookmarked entity

```sh
hatena-analyzer fetch-hatena-page-urls

hatena-analyzer fetch-bookmark

hatena-analyzer fetch-user-bm-count

hatena-analyzer view-timeseries

hatena-analyzer view-bookmark-details

hatena-analyzer view-summary
```

### use as Web Server

```sh
# Run as web server
hatena-analyzer web --port=8080

# request
curl http://localhost:8080/api/v1/fetch-page-url
```

## TODO

- [x] CLI Interface
- [x] Web Interface
- [ ] Web Handler response each
- [ ] analyze endpoint for `fetch-bookmark`, `view-bookmark-details`, `view-summary` at once
