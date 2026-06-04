ALTER TABLE songs
    ADD COLUMN artist_id text NOT NULL,
    ADD COLUMN album_id text NOT NULL,
    ADD COLUMN explicit_lyrics boolean DEFAULT false,
    ADD COLUMN isrc text NOT NULL,
    ADD COLUMN rank integer NOT NULL,
    ADD COLUMN link text,
    ADD COLUMN expicit_content_cover integer NOT NULL;
