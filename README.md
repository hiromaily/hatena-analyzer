# hatena-analyzer

Analyze bookmarked entity on [Hatena](https://b.hatena.ne.jp/hotentry/all)

## Requirements

- Golang 1.23+
- Docker

## Commands

```sh
# Help
go run ./cmd/analyzer/ -h
```

- `fetch-page-urls`: Fetched listed urls from Hatena page
- `fetch-bookmark`: Fetch bookmark entity information from url and save data to the database
- `fetch-user-bm-count`: Fetch user's bookmark count
- `view-timeseries`: View time series of bookmarked entity
- `view-bookmark-details`: View details of bookmarked entity
- `view-summary`: View summary of bookmarked entity

## TODO

- [x] CLI Interface
- [ ] Web Interface
