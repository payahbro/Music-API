package repository

import (
	"context"
	"database/sql"
)

type LikesRepository interface {
	CountLikes(ctx context.Context, id int64) (int64, error)
}

type LikesRepositoryImpl struct {
	Db *sql.DB
}

func NewLikeRepositoryImpl(db *sql.DB) LikesRepository {
	return &LikesRepositoryImpl{
		Db: db,
	}
}

func (l LikesRepositoryImpl) CountLikes(ctx context.Context, id int64) (int64, error) {
	var likes int64

	script := `
		SELECT COUNT(l.id_tracks) AS like_count
		FROM tracks t LEFT JOIN likes l ON t.id = l.id_tracks
		WHERE t.id = $1;
	`
	args := []any{id}
	row := l.Db.QueryRowContext(ctx, script, args...)

	err := row.Scan(&likes)
	if err != nil {
		return -1, err
	}
	return likes, err
}
