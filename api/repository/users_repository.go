package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"music-echo/api/domain/dao"
	"time"
)

type UsersRepository interface {
	Insert(ctx context.Context, users *dao.Users) error
	GetByEmail(ctx context.Context, email string) (*dao.Users, error)
	Update(ctx context.Context, users *dao.Users) error
	GetByToken(ctx context.Context, plainText, tokenScope string) (*dao.Users, error)
}

type UsersRepositoryImpl struct {
	Db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) UsersRepository {
	return UsersRepositoryImpl{
		Db: db,
	}
}

func (u UsersRepositoryImpl) Insert(ctx context.Context, users *dao.Users) error {
	script := `
		INSERT INTO users(name, email, password_hash)
    		VALUES ($1, $2, $3)
		RETURNING id, created_at, version;
	`
	args := []any{users.Name, users.Email, users.Password.Hash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := u.Db.QueryRowContext(ctx, script, args...)
	err := row.Scan(&users.Id, &users.CreatedAt, &users.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "email_unique"`:
			return errors.New("duplicate email")
		default:
			return err
		}
	}

	return nil
}

func (u UsersRepositoryImpl) GetByEmail(ctx context.Context, email string) (*dao.Users, error) {
	script := `
		SELECT id, created_at, name, email, password_hash, activated, version 
		FROM users
		WHERE email=$1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users dao.Users
	row := u.Db.QueryRowContext(ctx, script, email)
	err := row.Scan(
		&users.Id,
		&users.CreatedAt,
		&users.Name,
		&users.Email,
		&users.Password,
		&users.Activated,
		&users.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("record not found")
		default:
			return nil, err
		}
	}

	return &users, nil
}

func (u UsersRepositoryImpl) Update(ctx context.Context, users *dao.Users) error {
	script := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version;
	`
	args := []any{users.Name, users.Email, users.Password.Hash, users.Activated, users.Id, users.Version}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := u.Db.QueryRowContext(ctx, script, args...)
	err := row.Scan(&users.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "email_unique"`:
			return errors.New("duplicate email")
		case errors.Is(err, sql.ErrNoRows):
			return errors.New("edit conflict")
		default:
			return err
		}
	}

	return nil

}

func (u UsersRepositoryImpl) GetByToken(ctx context.Context, plainText, tokenScope string) (*dao.Users, error) {
	var user dao.Users
	hash := sha256.Sum256([]byte(plainText))

	script := `
	SELECT u.id, u.created_at, u.name, u.email, u.password_hash, u.activated, u.version
	FROM users u INNER JOIN token t ON u.id = t.user_id
	WHERE t.hash= $1 AND t.expiry > $2 AND t.scope=$3
	`

	args := []any{hash[:], time.Now(), tokenScope}

	err := u.Db.QueryRowContext(ctx, script, args...).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("no record")
		default:
			return nil, err
		}

	}

	return &user, nil
}
