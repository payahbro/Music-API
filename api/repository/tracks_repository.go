package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"music-echo/api/domain/dao"
	"music-echo/utils"
)

type TracksRepository interface {
	GetId(ctx context.Context, id int64) (*dao.Tracks, *dao.Artists, *int64, error)
	Insert(ctx context.Context, tracks *dao.Tracks) error
	Update(ctx context.Context, tracks *dao.Tracks) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, title, artist string, sorting utils.Sortings, paginating utils.Paginatings, genre []string) ([]*dao.Tracks, []*dao.Artists, []int64, int64, error)
}

type TracksRepositoryImpl struct {
	Db *sql.DB
}

func NewTracksRepositoryImpl(db *sql.DB) TracksRepository {
	return &TracksRepositoryImpl{Db: db}
}

func (t TracksRepositoryImpl) Insert(ctx context.Context, tracks *dao.Tracks) error {
	script := `
		INSERT INTO tracks(idartist, title, duration, year, genre) VALUES($1, $2, $3, $4, $5) RETURNING id, created_at
	`
	args := []any{tracks.IdArtist, tracks.Title, tracks.Duration, tracks.Year, pq.Array(tracks.Genre)}
	row := t.Db.QueryRowContext(ctx, script, args...)
	err := row.Scan(&tracks.Id, &tracks.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (t TracksRepositoryImpl) Update(ctx context.Context, tracks *dao.Tracks) error {
	script := `
		UPDATE tracks 
		SET idartist=$1, title=$2, duration=$3, year=$4, genre=$5, version=version+1 
		WHERE id=$6
		RETURNING version`

	args := []interface{}{tracks.IdArtist, tracks.Title, tracks.Duration, tracks.Year, pq.Array(tracks.Genre), tracks.Id}

	row := t.Db.QueryRowContext(ctx, script, args...)
	err := row.Scan(&tracks.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("track doesnt exist")
	}
	if err != nil {
		return err
	}

	return nil
}

func (t TracksRepositoryImpl) Delete(ctx context.Context, id int64) error {
	script := `
		DELETE
		FROM tracks
		WHERE id=$1;
	`
	row, err := t.Db.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("track doesnt exist")
	}

	return nil
}

func (t TracksRepositoryImpl) GetId(ctx context.Context, id int64) (*dao.Tracks, *dao.Artists, *int64, error) {
	script := `
		SELECT	t.id, t.created_at, t.idartist, t.title, t.duration, t.year, t.genre, t.version,
       			a.id AS artist_id, a.name AS artist_name,
       			COUNT(l.id_tracks) AS like_count 
		FROM tracks t 
		    LEFT JOIN likes l ON t.id = l.id_tracks
			LEFT JOIN artist a ON t.idartist = a.id
		WHERE t.id=$1
		GROUP BY t.id, a.id
    `
	args := []interface{}{id}
	row := t.Db.QueryRowContext(ctx, script, args...)

	var track dao.Tracks
	var artist dao.Artists
	var likes int64
	err := row.Scan(
		&track.Id,
		&track.CreatedAt,
		&track.IdArtist,
		&track.Title,
		&track.Duration,
		&track.Year,
		pq.Array(&track.Genre),
		&track.Version,
		&artist.Id,
		&artist.Name,
		&likes)

	if err != nil {
		return nil, nil, nil, err
	}
	return &track, &artist, &likes, nil

}

func (t TracksRepositoryImpl) GetAll(ctx context.Context, title, artist string, sorting utils.Sortings, paginating utils.Paginatings, genre []string) ([]*dao.Tracks, []*dao.Artists, []int64, int64, error) {
	var script = fmt.Sprintf(`
		SELECT 	COUNT(*) OVER(), t.id, t.created_at, t.idartist, t.title, t.duration, t.year, t.genre, t.version,
          		a.id AS artist_id, a.name AS artist_name,
          		COUNT(l.id_tracks) AS likes_count
		FROM tracks t
         		LEFT JOIN artist a ON t.idartist = a.id
         		LEFT JOIN likes l ON t.id = l.id_tracks
		WHERE (to_tsvector('simple', t.title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
		  		AND (to_tsvector('simple', a.name) @@ plainto_tsquery('simple', $2) OR $2 = '') 
		  		AND(t.genre @> $3 OR $3='{}')
		GROUP BY t.id, a.id
		ORDER BY %s %s, t.id ASC
		LIMIT $4 OFFSET $5`, sorting.SortName(), sorting.SortDirection())

	var args = []interface{}{title, artist, pq.Array(genre), paginating.Limit(), paginating.Offset()}
	var rows, err = t.Db.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	var tracks []*dao.Tracks
	var artists []*dao.Artists
	var likes []int64
	var totalRecords int64
	for rows.Next() {
		var track dao.Tracks
		var artist dao.Artists
		var like int64
		err = rows.Scan(
			&totalRecords,
			&track.Id,
			&track.CreatedAt,
			&track.IdArtist,
			&track.Title,
			&track.Duration,
			&track.Year,
			pq.Array(&track.Genre),
			&track.Version,
			&artist.Id,
			&artist.Name,
			&like,
		)
		if err != nil {
			return nil, nil, nil, 0, err
		}

		tracks = append(tracks, &track)
		artists = append(artists, &artist)
		likes = append(likes, like)
	}

	return tracks, artists, likes, totalRecords, nil

}
