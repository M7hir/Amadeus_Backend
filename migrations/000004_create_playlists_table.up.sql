CREATE TABLE playlists (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    name text NOT NULL,
    description text,
    is_public boolean NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1,
    CONSTRAINT playlists_name_check CHECK (length(name) > 0)
);

CREATE TRIGGER playlists_touch_updated_at
BEFORE UPDATE ON playlists
FOR EACH ROW
EXECUTE FUNCTION touch_updated_at();

CREATE INDEX idx_playlists_user_id ON playlists(user_id);
