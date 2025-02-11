# hatena-fake-detector

例えば、以下のA氏の記事は不自然なブックマーク数の上昇や、不自然なポジティブコメントで溢れている。
またネガティブなコメントはすぐに非表示になるため、高確率で不正をしていると思われる。

## Requirements

- Golang
- Docker

## Commands

- `fetch-bookmark`: Fetch bookmarked entity from URL and save data to the database
- `view-summary`: View time series data of the summary of bookmarked entity
- `fetch-user-info`:
- `view-page-score`:

## 必要な機能

- 記事に対して
  - 記事のブックマーク数の取得
  - 記事に対して、ブックマークしているユーザーの取得
- ユーザーに対して
  - いつ作成されたユーザーか
  - 過去のブックマーク数の取得 (新規ユーザーであるかどうか)
  - そのユーザーが他の対象のURLに対してもブックマークしているかを確認
- つまり、すべてのURLに対しての全ユーザー情報として、過去のブックマーク数、A氏サイトのブックマーク数、この２つの情報が出せれば良い
- TODO: 不正をスコア化する
