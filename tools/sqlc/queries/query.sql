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

-- name: InsertURL :execlastid
-- @desc: insert url if not existed and return url_id
INSERT INTO
  URLs (url_address)
VALUES
  ($1)
ON CONFLICT (url_address) DO NOTHING
RETURNING
  url_id;

-- name: InsertUser :execlastid
-- @desc: insert user if not existed and return user_id
INSERT INTO
  Users (user_name)
VALUES
  ($1)
ON CONFLICT (user_name) DO NOTHING
RETURNING
  user_id;
