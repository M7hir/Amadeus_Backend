ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

CREATE TABLE user_identities (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider text NOT NULL,     -- e.g., 'google', 'facebook'
    provider_id text NOT NULL,  -- The unique ID returned by the provider
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    
    -- Ensure a specific provider ID can only be linked to one user
    CONSTRAINT user_identities_provider_unique UNIQUE (provider, provider_id),
    CONSTRAINT user_identities_provider_check CHECK (length(provider) > 0),
    CONSTRAINT user_identities_provider_id_check CHECK (length(provider_id) > 0)
);

CREATE TRIGGER user_identities_touch_updated_at
BEFORE UPDATE ON user_identities
FOR EACH ROW
EXECUTE FUNCTION touch_updated_at();

CREATE INDEX idx_user_identities_user_id ON user_identities(user_id);
CREATE INDEX idx_user_identities_provider_id ON user_identities(provider, provider_id);
