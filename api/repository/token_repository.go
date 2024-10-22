package repository

import (
	"context"
	"database/sql"
	"music-echo/api/domain/dao"
	"time"
)

type TokenRepository interface {
	Insert(ctx context.Context, tokens *dao.Token) error
	Delete(ctx context.Context, userId int64, scope string) error
}

type TokenRepositoryImpl struct {
	Db *sql.DB
}

func NewTokenRepositoryImpl(db *sql.DB) TokenRepository {
	return &TokenRepositoryImpl{
		Db: db,
	}
}

func (t TokenRepositoryImpl) Insert(ctx context.Context, tokens *dao.Token) error {
	script := `
		INSERT INTO token(hash, user_id, expiry, scope)
    	VALUES($1, $2, $3, $4)
	`
	args := []any{tokens.Hash, tokens.UserId, tokens.Expiry, tokens.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.Db.ExecContext(ctx, script, args...)
	if err != nil {
		return err
	}

	return nil
}

func (t TokenRepositoryImpl) Delete(ctx context.Context, userId int64, scope string) error {
	script := `
		DELETE
		FROM token
		WHERE user_id = $1 AND scope = $2
	`
	args := []any{userId, scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.Db.ExecContext(ctx, script, args...)
	if err != nil {
		return err
	}

	return nil

}
