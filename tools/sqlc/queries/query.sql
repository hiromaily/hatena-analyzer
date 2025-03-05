-- name: GetURLsByURLAddresses :many
-- @desc: get url information by url address
SELECT DISTINCT ON (u.url_address)
  u.url_id, u.url_address, u.category_code, u.title, u.bookmark_count, u.named_user_count, u.private_user_rate
FROM
  URLs u
WHERE
  u.is_deleted = FALSE
AND
  u.url_address = ANY($1::text[])
ORDER BY
  u.url_address, u.url_id;

-- name: GetAllURLs :many
-- @desc: get all url addresses
SELECT DISTINCT ON (u.url_address)
  u.url_id, u.url_address, u.category_code, u.title, u.bookmark_count, u.named_user_count, u.private_user_rate
FROM
  URLs u
WHERE
  u.is_deleted = FALSE
ORDER BY
  u.url_address, u.url_id;

-- name: GetUrlID :one
-- @desc: get target url_id by url address
SELECT
  u.url_id
FROM
  URLs u
WHERE
  u.url_address = $1;

-- name: GetURLsByPrivateRate :many
-- @desc: get urls by private_user_rate
SELECT
  u.url_id, u.url_address, u.category_code, u.title, u.bookmark_count, u.named_user_count, u.private_user_rate
FROM 
  URLs u
WHERE 
  private_user_rate >= $1 
ORDER BY 
  private_user_rate DESC;


-- name: GetUserNames :many
-- @desc: get users
SELECT
  u.user_name
FROM
  Users u
WHERE
  u.is_deleted = FALSE;

-- name: GetUserNamesByURL :many
-- @desc: get target users by url
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
  AND url.url_address = $1;

-- name: GetUserNamesByURLs :many
-- @desc: get target users by multiple urls
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
  AND url.url_address = ANY($1::text[]);

-- name: GetUsersByURL :many
-- @desc: get target users by url
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
  u.bookmark_count DESC;

-- name: UpdateUserBookmarkCount :one
-- @desc: update user bookmark count and return url_id
UPDATE Users
  SET bookmark_count = $1, updated_at = CURRENT_TIMESTAMP
WHERE user_name = $2 
RETURNING
  user_id;

-- name: InsertURL :one
-- @desc: Deprecated. insert url if not existed and return url_id
WITH insert_result AS (
	INSERT INTO URLs (url_address, category_code)
	VALUES ($1, $2)
	ON CONFLICT (url_address, category_code) DO NOTHING
	RETURNING url_id
)
SELECT url_id FROM insert_result
UNION ALL
SELECT url_id FROM URLs WHERE url_address = $1 AND category_code = $2 LIMIT 1;

-- name: InsertURLs :copyfrom
-- @desc: Deprecated. insert urls if not existed
INSERT INTO
  URLs (url_address, category_code)
VALUES
  ($1, $2);

-- name: BulkInsertUrls :exec
-- @desc: insert urls by stored procedure. conflicts must be ignored. arg1: array of urls, arg2: array of category, arg3: array of isAll flag.
-- name: BulkInsertURLs :exec
CALL bulk_insert_urls($1, $2, $3);

-- name: UpsertURL :one
-- @desc: insert url if not existed, update url with is_deleted=false if existed
INSERT INTO URLs (url_address, title, bookmark_count, named_user_count, private_user_rate) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (url_address) 
DO UPDATE SET
    bookmark_count = $3,
    named_user_count = $4,
    private_user_rate = $5,
    is_deleted = FALSE,
    updated_at = EXCLUDED.updated_at 
RETURNING url_id;

-- name: UpdateURL :execrows
-- @desc: update url with bookmark_count, named_user_count, private_user_rate
UPDATE URLs
SET
    title = $1,
    bookmark_count = $2,
    named_user_count = $3,
    private_user_rate = $4
WHERE
    url_id = $5;

-- name: InsertUser :one
-- @desc: Deprecated!!! insert user if not existed and return user_id
INSERT INTO
  Users (user_name)
VALUES
  ($1)
ON CONFLICT (user_name) DO NOTHING
RETURNING
  user_id;

-- name: UpsertUser :one
-- @desc: insert user if not existed, update user with is_deleted=false if existed
INSERT INTO Users (user_name) 
VALUES ($1)
ON CONFLICT (user_name) 
DO UPDATE SET 
    is_deleted = FALSE,
    updated_at = EXCLUDED.updated_at 
RETURNING user_id;

-- name: UpsertUserURLs :exec
-- @desc: insert UserURLs if not existed, update UserURLs with is_deleted=false if existed
INSERT INTO UserURLs (user_id, url_id) 
VALUES ($1, $2)
ON CONFLICT (user_id, url_id) 
DO UPDATE SET 
    is_deleted = FALSE,
    updated_at = EXCLUDED.updated_at;

-- name: GetBookmarkedUsersURLCounts :many
-- @desc: Not used. Count each user's bookmarked urls
SELECT 
    user_id, COUNT(user_id) AS url_count
FROM 
    UserURLs
WHERE
    url_id in (1,2,3,4)
GROUP BY 
    user_id
ORDER BY 
    url_count DESC;

-- name: CountGetBookmarkedUsersURLCounts :one
-- @desc: Not used. Count target that each user's bookmarked urls
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
  ) AS subquery;

-- name: GetAveragePrivateUserRates :many
-- @desc: get average private user rates on all categories
SELECT
  category_code, AVG(private_user_rate) AS average_private_user_rate
FROM 
  URLs
WHERE 
  is_deleted = FALSE
GROUP BY 
  category_code;
