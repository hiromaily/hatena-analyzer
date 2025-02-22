// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sqlcgen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countGetBookmarkedUsersURLCounts = `-- name: CountGetBookmarkedUsersURLCounts :one
SELECT
  COUNT(*)
FROM
  (
    SELECT
      user_id,
      COUNT(user_id) AS url_count
    FROM
      UserURLs
    WHERE
      url_id IN (1, 2, 3, 4)
    GROUP BY
      user_id
    HAVING
      COUNT(user_id) = 4
  ) AS subquery
`

// @desc: Count target that each user's bookmarked urls
func (q *Queries) CountGetBookmarkedUsersURLCounts(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countGetBookmarkedUsersURLCounts)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getAllURLs = `-- name: GetAllURLs :many
SELECT
  u.url_id, u.url_address
FROM
  URLs u
WHERE
  u.is_deleted = FALSE
`

type GetAllURLsRow struct {
	UrlID      int32
	UrlAddress string
}

// @desc: get all url addresses
func (q *Queries) GetAllURLs(ctx context.Context) ([]GetAllURLsRow, error) {
	rows, err := q.db.Query(ctx, getAllURLs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllURLsRow
	for rows.Next() {
		var i GetAllURLsRow
		if err := rows.Scan(&i.UrlID, &i.UrlAddress); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBookmarkedUsersURLCounts = `-- name: GetBookmarkedUsersURLCounts :many
SELECT 
    user_id, COUNT(user_id) AS url_count
FROM 
    UserURLs
WHERE
    url_id in (1,2,3,4)
GROUP BY 
    user_id
ORDER BY 
    url_count DESC
`

type GetBookmarkedUsersURLCountsRow struct {
	UserID   int32
	UrlCount int64
}

// @desc: Count each user's bookmarked urls
func (q *Queries) GetBookmarkedUsersURLCounts(ctx context.Context) ([]GetBookmarkedUsersURLCountsRow, error) {
	rows, err := q.db.Query(ctx, getBookmarkedUsersURLCounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBookmarkedUsersURLCountsRow
	for rows.Next() {
		var i GetBookmarkedUsersURLCountsRow
		if err := rows.Scan(&i.UserID, &i.UrlCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUrlID = `-- name: GetUrlID :one
SELECT
  u.url_id
FROM
  URLs u
WHERE
  u.url_address = $1
`

// @desc: get target url_id by url address
func (q *Queries) GetUrlID(ctx context.Context, urlAddress string) (int32, error) {
	row := q.db.QueryRow(ctx, getUrlID, urlAddress)
	var url_id int32
	err := row.Scan(&url_id)
	return url_id, err
}

const getUserNames = `-- name: GetUserNames :many
SELECT
  u.user_name
FROM
  Users u
WHERE
  u.is_deleted = FALSE
`

// @desc: get users
func (q *Queries) GetUserNames(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getUserNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var user_name string
		if err := rows.Scan(&user_name); err != nil {
			return nil, err
		}
		items = append(items, user_name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserNamesByURL = `-- name: GetUserNamesByURL :many
SELECT
  u.user_name
FROM
  Users u
  INNER JOIN UserURLs uu ON u.user_id = uu.user_id
  INNER JOIN URLs url ON uu.url_id = url.url_id
WHERE
  u.is_deleted = FALSE
  AND url.is_deleted = FALSE
  AND uu.is_deleted = FALSE
  AND url.url_address = $1
`

// @desc: get target users by url
func (q *Queries) GetUserNamesByURL(ctx context.Context, urlAddress string) ([]string, error) {
	rows, err := q.db.Query(ctx, getUserNamesByURL, urlAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var user_name string
		if err := rows.Scan(&user_name); err != nil {
			return nil, err
		}
		items = append(items, user_name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserNamesByURLs = `-- name: GetUserNamesByURLs :many
SELECT
  u.user_name
FROM
  Users u
  INNER JOIN UserURLs uu ON u.user_id = uu.user_id
  INNER JOIN URLs url ON uu.url_id = url.url_id
WHERE
  u.is_deleted = FALSE
  AND url.is_deleted = FALSE
  AND uu.is_deleted = FALSE
  AND url.url_address = ANY($1::text[])
`

// @desc: get target users by multiple urls
func (q *Queries) GetUserNamesByURLs(ctx context.Context, dollar_1 []string) ([]string, error) {
	rows, err := q.db.Query(ctx, getUserNamesByURLs, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var user_name string
		if err := rows.Scan(&user_name); err != nil {
			return nil, err
		}
		items = append(items, user_name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersByURL = `-- name: GetUsersByURL :many
SELECT
  u.user_name, u.bookmark_count
FROM
  Users u
  INNER JOIN UserURLs uu ON u.user_id = uu.user_id
  INNER JOIN URLs url ON uu.url_id = url.url_id
WHERE
  u.is_deleted = FALSE
  AND url.is_deleted = FALSE
  AND uu.is_deleted = FALSE
  AND url.url_address = $1
ORDER BY
  u.bookmark_count DESC
`

type GetUsersByURLRow struct {
	UserName      string
	BookmarkCount pgtype.Int4
}

// @desc: get target users by url
func (q *Queries) GetUsersByURL(ctx context.Context, urlAddress string) ([]GetUsersByURLRow, error) {
	rows, err := q.db.Query(ctx, getUsersByURL, urlAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersByURLRow
	for rows.Next() {
		var i GetUsersByURLRow
		if err := rows.Scan(&i.UserName, &i.BookmarkCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertURL = `-- name: InsertURL :one
WITH insert_result AS (
	INSERT INTO URLs (url_address, category_code)
	VALUES ($1, $2)
	ON CONFLICT (url_address, category_code) DO NOTHING
	RETURNING url_id
)
SELECT url_id FROM insert_result
UNION ALL
SELECT url_id FROM URLs WHERE url_address = $1 LIMIT 1
`

type InsertURLParams struct {
	UrlAddress   string
	CategoryCode pgtype.Text
}

// @desc: insert url if not existed and return url_id
func (q *Queries) InsertURL(ctx context.Context, arg InsertURLParams) (int32, error) {
	row := q.db.QueryRow(ctx, insertURL, arg.UrlAddress, arg.CategoryCode)
	var url_id int32
	err := row.Scan(&url_id)
	return url_id, err
}

type InsertURLsParams struct {
	UrlAddress   string
	CategoryCode pgtype.Text
}

const insertUser = `-- name: InsertUser :one
INSERT INTO
  Users (user_name)
VALUES
  ($1)
ON CONFLICT (user_name) DO NOTHING
RETURNING
  user_id
`

// @desc: Deprecated!!! insert user if not existed and return user_id
func (q *Queries) InsertUser(ctx context.Context, userName string) (int32, error) {
	row := q.db.QueryRow(ctx, insertUser, userName)
	var user_id int32
	err := row.Scan(&user_id)
	return user_id, err
}

const updateUserBookmarkCount = `-- name: UpdateUserBookmarkCount :one
UPDATE Users
  SET bookmark_count = $1, updated_at = CURRENT_TIMESTAMP
WHERE user_name = $2 
RETURNING
  user_id
`

type UpdateUserBookmarkCountParams struct {
	BookmarkCount pgtype.Int4
	UserName      string
}

// @desc: update user bookmark count and return url_id
func (q *Queries) UpdateUserBookmarkCount(ctx context.Context, arg UpdateUserBookmarkCountParams) (int32, error) {
	row := q.db.QueryRow(ctx, updateUserBookmarkCount, arg.BookmarkCount, arg.UserName)
	var user_id int32
	err := row.Scan(&user_id)
	return user_id, err
}

const upsertUser = `-- name: UpsertUser :one
INSERT INTO Users (user_name) 
VALUES ($1)
ON CONFLICT (user_name) 
DO UPDATE SET 
    is_deleted = FALSE,
    updated_at = EXCLUDED.updated_at 
RETURNING user_id
`

// @desc: insert user if not existed, update user with is_deleted=false if existed
func (q *Queries) UpsertUser(ctx context.Context, userName string) (int32, error) {
	row := q.db.QueryRow(ctx, upsertUser, userName)
	var user_id int32
	err := row.Scan(&user_id)
	return user_id, err
}

const upsertUserURLs = `-- name: UpsertUserURLs :exec
INSERT INTO UserURLs (user_id, url_id) 
VALUES ($1, $2)
ON CONFLICT (user_id, url_id) 
DO UPDATE SET 
    is_deleted = FALSE,
    updated_at = EXCLUDED.updated_at
`

type UpsertUserURLsParams struct {
	UserID int32
	UrlID  int32
}

// @desc: insert UserURLs if not existed, update UserURLs with is_deleted=false if existed
func (q *Queries) UpsertUserURLs(ctx context.Context, arg UpsertUserURLsParams) error {
	_, err := q.db.Exec(ctx, upsertUserURLs, arg.UserID, arg.UrlID)
	return err
}
