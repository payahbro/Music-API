ALTER TABLE likes
    ADD CONSTRAINT unique_user_tracks UNIQUE (id_users, id_tracks);