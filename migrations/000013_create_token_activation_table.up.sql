CREATE TABLE token(
    hash BYTEA PRIMARY KEY,
    user_id BIGINT NOT NULL,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL,
    CONSTRAINT fk_token_user FOREIGN KEY(user_id) REFERENCES users ON DELETE CASCADE
)