# hatena-fake-detector

例えば、以下のA氏の記事は不自然なブックマーク数の上昇や、不自然なポジティブコメントで溢れている。
またネガティブなコメントはすぐに非表示になるため、高確率で不正をしていると思われる。

https://note.com/simplearchitect

## 前提条件

A氏のサイトのURLは事前にプログラムで登録しておく。もしくは簡単に追加できるようにしておく。

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

## 利用するストレージ

### InfluxDB

- bucket: Bookmark
- measurement(table): url
- tag
  - title: URLのタイトル
- field
  - count: ブックマーク数
  - user_num: ユーザー数

### MongoDB





## [はてなブックマーク REST API](https://developer.hatena.ne.jp/ja/documents/bookmark/apis/rest/)

### [はてなブックマーク エントリー情報取得API](https://developer.hatena.ne.jp/ja/documents/bookmark/apis/getinfo)

- URL: `https://b.hatena.ne.jp/entry/json/任意のURL`
- 例: https://b.hatena.ne.jp/entry/json/https://note.com/simplearchitect/n/nadc0bcdd5b3d
- Response: JSON

curl例

```sh
curl -s https://b.hatena.ne.jp/entry/json/https://note.com/simplearchitect/n/nadc0bcdd5b3d | jq '{title: .title, count: .count}'
curl -s https://b.hatena.ne.jp/entry/json/https://note.com/simplearchitect/n/nadc0bcdd5b3d | jq '.bookmarks[].user'
curl -s https://b.hatena.ne.jp/entry/json/https://note.com/simplearchitect/n/nadc0bcdd5b3d | jq '.bookmarks | length'
```

https://b.hatena.ne.jp/json/osugi3y/

#### count数とブックマーク数が異なるのはなぜか？

古いブックマーク情報や削除されたブックマーク情報は含まれていない可能性がある。

### [はてなブックマーク ユーザー情報 API](https://developer.hatena.ne.jp/ja/documents/bookmark/apis/rest/my/)

認証したユーザーの情報を取得 (つまり自分だけか)

- URL: https://bookmark.hatenaapis.com/rest/1/osugi3y

## 生成AIへの頼み方

今から説明することをgolangでプログラムを書いてください。

まず、以下にURLがあります。

```
https://note.com/simplearchitect/n/nadc0bcdd5b3d
https://note.com/simplearchitect/n/n871f29ffbfac
https://note.com/simplearchitect/n/n86a95bc19b4c
https://note.com/simplearchitect/n/nfd147540e3db
```

これらは、はてなブックマークというサービス上でも登録されている状態です。
そして、はてなブックマークのAPIを使って、各URLの情報をJSONで取得することが可能です。

`https://b.hatena.ne.jp/entry/json/任意のURL` で取得することができます。
例えば最初に教えたURLの1つを例にすると、
`https://b.hatena.ne.jp/entry/json/https://note.com/simplearchitect/n/nadc0bcdd5b3d` という形で情報をJSON形式で取得することができます。

JSONは以下のような構造になっています。

```json
{
    "url": "https://example.com/news-0001",
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
    "requested_url": "https://example.com/news-0001"
}
```

ここから取得したい情報は、タイトル、ブックマークカウント数、ユーザー一覧、ユーザー数、各ユーザーのコメントの有無です。
仮にjqコマンドを利用するのであれば、以下のように取得可能です。

```sh
curl -s https://b.hatena.ne.jp/entry/json/https://note.com/simplearchitect/n/nadc0bcdd5b3d | jq '{title: .title, count: .count}'
curl -s https://b.hatena.ne.jp/entry/json/https://note.com/simplearchitect/n/nadc0bcdd5b3d | jq '.bookmarks[].user'
curl -s https://b.hatena.ne.jp/entry/json/https://note.com/simplearchitect/n/nadc0bcdd5b3d | jq '.bookmarks | length'
```

最終的に期待する構造体は以下の通りですが、
JSONから取得できるCountと、User配列の要素数はなぜか異なるため、User配列の要素数はUser Countという項目で、別途表示したいです。

```go
type Bookmark struct {
   title string 
   count int
   users map[string]User // keyはuserの名前をセットします
}

type User struct {
   name string
   isCommented bool
   isDeleted bool
}
```

最後に複数のURLのuserをマージして、その総数を知りたいです。
以下のようなイメージです。

```go
type BookmarkUsers struct {
   users map[string]User
}

type User struct {
   name string
   totalBookmarkNum int
}
```
