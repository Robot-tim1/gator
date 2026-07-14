-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)

SELECT
    iff.*,
    feeds.name as feed_name,
    users.name as user_name
FROM inserted_feed_follow AS iff
INNER JOIN users
    ON iff.user_id = users.id
INNER JOIN feeds
    ON iff.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT 
    users.name as user_name,
    feeds.name as feed_name, 
    feeds.url as feed_url,
    feed_follows.*
FROM 
    feed_follows
INNER JOIN users
    ON users.id = feed_follows.user_id
INNER JOIN feeds
    on feeds.id = feed_follows.feed_id
WHERE users.id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
where user_id = $1 and feed_id = $2;