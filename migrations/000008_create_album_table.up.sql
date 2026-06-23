CREATE TABLE album (
    id text PRIMARY KEY,
    name text NOT NULL,
    cover text,
    cover_big text,
    cover_medium text,
    cover_small text,
    cover_xl text,
    link text,
    tracklist text,
    genre_id int NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT album_name_check CHECK (length(name) > 0),
    CONSTRAINT album_genre_id_check CHECK (genre_id >= 0)
);


CREATE TRIGGER album_touch_updated_at
BEFORE UPDATE ON album
FOR EACH ROW
EXECUTE FUNCTION touch_updated_at();

CREATE INDEX idx_album_name ON album(name);