CREATE TABLE playlist_songs (
    playlist_id uuid NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    song_id text NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    added_at timestamptz NOT NULL DEFAULT now(),
    position numeric(20,10) NOT NULL,
    PRIMARY KEY (playlist_id, song_id),
    CONSTRAINT playlist_songs_position_check CHECK (position > 0)
);

CREATE INDEX idx_playlist_songs_song_id ON playlist_songs(song_id);
CREATE INDEX idx_playlist_songs_playlist_id_position ON playlist_songs(playlist_id, position);
CREATE INDEX idx_playlist_songs_playlist_id_added_at ON playlist_songs(playlist_id, added_at);
