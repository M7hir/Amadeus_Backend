CREATE OR REPLACE FUNCTION touch_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    first_name text NOT NULL,
    last_name text NOT NULL,
    email citext NOT NULL UNIQUE,
    password_hash text NOT NULL,
    activated boolean NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1,
    CONSTRAINT users_first_name_check CHECK (length(first_name) > 0),
    CONSTRAINT users_last_name_check CHECK (length(last_name) > 0),
    CONSTRAINT users_email_check CHECK (position('@' in email) > 1)
);

CREATE TRIGGER users_touch_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION touch_updated_at();

CREATE INDEX idx_users_activated ON users(activated);
