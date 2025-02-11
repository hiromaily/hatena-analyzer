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

-- name: BookmarkedUsersCounts :many
-- @desc: Count each user's bookmarked urls
WITH given_urls AS (
    SELECT UNNEST(ARRAY[1, 2, 3, 4]) AS url_id
)
SELECT 
    U.user_id,
    U.user_name,
    COUNT(UR.url_id) as url_count
FROM 
    Users U
JOIN 
    UserURLs UR ON U.user_id = UR.user_id
JOIN 
    URLs ON UR.url_id = URLs.url_id
JOIN 
    given_urls GU ON UR.url_id = GU.url_id
WHERE 
    U.is_deleted = FALSE
AND 
    URLs.is_deleted = FALSE
AND 
    UR.is_deleted = FALSE
GROUP BY 
    U.user_id, U.user_name
ORDER BY 
    url_count DESC;
