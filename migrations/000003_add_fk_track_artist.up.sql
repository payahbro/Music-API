ALTER TABLE tracks
    ADD CONSTRAINT fk_track_artist FOREIGN KEY (idartist) REFERENCES artist(id);
