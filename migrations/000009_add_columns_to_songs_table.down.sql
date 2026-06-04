ALTER TABLE songs
    DROP COLUMN IF EXISTS artist_id,
    DROP COLUMN IF EXISTS album_id,
    DROP COLUMN IF EXISTS expicit_content_cover,
    DROP COLUMN IF EXISTS rank,
    DROP COLUMN IF EXISTS isrc,
    DROP COLUMN IF EXISTS link,
    DROP COLUMN IF EXISTS explicit_lyrics;