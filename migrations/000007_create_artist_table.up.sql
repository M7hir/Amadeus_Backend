CREATE TABLE artist (
    id text PRIMARY KEY,
    picture text,
    picture_big text,
    picture_small text,
    picture_medium text,
    picture_xl text,
    tracklist text,
    link text,
    title text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT artist_title_check CHECK (length(title) > 0)
);

CREATE TRIGGER artist_touch_updated_at
BEFORE UPDATE ON artist
FOR EACH ROW
EXECUTE FUNCTION touch_updated_at();

CREATE INDEX idx_artist_title ON artist(title);