CREATE TABLE songs (
    id text PRIMARY KEY,
    title text NOT NULL,
    artist_name text NOT NULL,
    duration_ms integer NOT NULL,
    cover_image_url text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT songs_title_check CHECK (length(title) > 0),
    CONSTRAINT songs_artist_name_check CHECK (length(artist_name) > 0),
    CONSTRAINT songs_duration_ms_check CHECK (duration_ms > 0)
);

CREATE TRIGGER songs_touch_updated_at
BEFORE UPDATE ON songs
FOR EACH ROW
EXECUTE FUNCTION touch_updated_at();

CREATE INDEX idx_songs_artist_name ON songs(artist_name);
