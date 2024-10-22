CREATE TABLE IF NOT EXISTS likes(
    id_users INTEGER,
    id_tracks INTEGER,
    CONSTRAINT fk_user_likes FOREIGN KEY(id_users) REFERENCES users(id),
    CONSTRAINT fk_tracks_likes FOREIGN KEY (id_tracks) REFERENCES tracks(id)
)