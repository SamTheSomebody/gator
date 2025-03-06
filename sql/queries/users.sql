-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1
LIMIT 1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: Reset :exec
DELETE FROM users;

-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
  VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
  )
  RETURNING *
) SELECT
  inserted_feed_follow.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds
ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users
ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, users.name, feeds.name FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: RemoveFeedFollow :exec
DELETE FROM feed_follows USING feeds
WHERE feed_follows.feed_id = feeds.id
AND feed_follows.user_id = $1
AND feeds.url = $2;
