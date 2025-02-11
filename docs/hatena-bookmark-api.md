# [はてなブックマーク REST API](https://developer.hatena.ne.jp/ja/documents/bookmark/apis/rest/)

## [はてなブックマーク エントリー情報取得API](https://developer.hatena.ne.jp/ja/documents/bookmark/apis/getinfo)

- URL: `https://b.hatena.ne.jp/entry/json/任意のURL`
- 例: https://b.hatena.ne.jp/entry/json/https://www.google.co.jp/
- Response: JSON

### 取得できるレスポンス

```json
{
    "url": "https://www.google.co.jp/",
    "eid": "4765841514130150390",
    "bookmarks": [
        {
            "user": "foo",
            "timestamp": "2025/02/05 21:32",
            "comment": "foo foo foo"
        },
        {
            "user": "bar",
            "timestamp": "2025/02/05 22:10",
            "comment": "bar bar bar"
        }
    ],
    "title": "\u30c6\u30b9\u30c8\u30c6\u30b9\u30c8\u30c6\u30b9\u30c8",
    "count": 123,
    "requested_url": "https://www.google.co.jp/"
}
```

### curl例

```sh
curl -s https://b.hatena.ne.jp/entry/json/https://www.google.co.jp/ | jq '{title: .title, count: .count}'
curl -s https://b.hatena.ne.jp/entry/json/https://www.google.co.jp/ | jq '.bookmarks[].user'
curl -s https://b.hatena.ne.jp/entry/json/https://www.google.co.jp/ | jq '.bookmarks | length'
```

### count数とブックマーク配列の要素数が異なるのはなぜか？

ブックマーク後にそのアカウントを削除したケースがこれに該当すると思われる。
つまり、この差分が大きいページのブックマーク数は不正の可能性が高い。

## [はてなブックマーク ユーザー情報 API](https://developer.hatena.ne.jp/ja/documents/bookmark/apis/rest/my/)

認証したユーザーの情報を取得 (つまり自分だけか)

