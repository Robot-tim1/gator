package database

import (
	"context"
	"fmt"
	"strings"
)

// could just be a join query or something, but I don't care I wanted to do a dynamic query
func (q *Queries) GetNextFeedFromFollows(ctx context.Context, feeds []GetFeedFollowsForUserRow) (Feed, error) {
	placeholders := make([]string, len(feeds))
	args := make([]any, len(feeds))

	for i, feed := range feeds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = feed.FeedID
	}

	query := fmt.Sprintf(`
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds
WHERE id IN (%s)
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1`, strings.Join(placeholders, ","))

	row := q.db.QueryRowContext(ctx, query, args...)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}
