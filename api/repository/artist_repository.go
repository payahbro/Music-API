package repository

import (
	"context"
	"database/sql"
	"music-echo/api/domain/dao"
)

type ArtistRepository interface {
	GetByName(ctx context.Context, name string) (*dao.Artists, error)
	GetById(ctx context.Context, id int64) (*dao.Artists, error)
}

type ArtistRepositoryImpl struct {
	Db *sql.DB
}

func NewArtistRepositoryImpl(db *sql.DB) ArtistRepository {
	return &ArtistRepositoryImpl{Db: db}
}

func (a ArtistRepositoryImpl) GetByName(ctx context.Context, name string) (*dao.Artists, error) {
	var artist dao.Artists
	script := "SELECT * FROM artist WHERE name=$1"
	args := []any{name}
	row := a.Db.QueryRowContext(ctx, script, args...)
	err := row.Scan(&artist.Id, &artist.Name)
	if err != nil {
		return nil, err
	}
	return &artist, err

}

func (a ArtistRepositoryImpl) GetById(ctx context.Context, id int64) (*dao.Artists, error) {
	var artist dao.Artists
	script := "SELECT * FROM artist WHERE id=$1"
	args := []any{id}
	row := a.Db.QueryRowContext(ctx, script, args...)
	err := row.Scan(&artist.Id, &artist.Name)
	if err != nil {
		return nil, err
	}

	return &artist, nil
}
