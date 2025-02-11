-- name: GetUsersByURL :many
-- @desc: get target users by url
SELECT u.*
FROM Users u
INNER JOIN UserURLs uu ON u.user_id = uu.user_id
INNER JOIN URLs url ON uu.url_id = url.url_id
WHERE url.url_address = $1;
