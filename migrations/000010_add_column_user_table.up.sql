ALTER TABLE users
    ADD COLUMN email CITEXT NOT NULL DEFAULT 'testing@gmail.com',
    ADD COLUMN password_hash BYTEA NOT NULL DEFAULT 'h4ny4t35t1ng',
    ADD COLUMN activated BOOL NOT NULL DEFAULT FALSE,
    ADD COLUMN version INTEGER NOT NULL DEFAULT 1


