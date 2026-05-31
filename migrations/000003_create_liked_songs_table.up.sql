CREATE TABLE liked_songs (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    song_id text NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    liked_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, song_id)
);

CREATE INDEX idx_liked_songs_song_id ON liked_songs(song_id);
