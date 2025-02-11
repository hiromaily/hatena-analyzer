-- name: GetUrlID :one
-- @desc: get target url_id by url address
SELECT
  u.url_id
FROM
  URLs u
WHERE
  u.url_address = $1;

-- name: GetUsersByURL :many
-- @desc: get target users by url
SELECT
  u.*
FROM
  Users u
  INNER JOIN UserURLs uu ON u.user_id = uu.user_id
  INNER JOIN URLs url ON uu.url_id = url.url_id
WHERE
  url.url_address = $1;

-- name: InsertURL :one
-- @desc: insert url if not existed and return url_id
INSERT INTO
  URLs (url_address)
VALUES
  ($1)
ON CONFLICT (url_address) DO NOTHING
RETURNING
  url_id;

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
-- @desc: Count each user's bookmarked urls
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
-- @desc: Count target that each user's bookmarked urls
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
