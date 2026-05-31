DROP TRIGGER IF EXISTS user_identities_touch_updated_at ON user_identities;
DROP TABLE IF EXISTS user_identities;

-- Note: Re-adding NOT NULL to password_hash is commented out because
-- it will fail if any users exist with NULL password_hash (OAuth users).
-- In a real rollback scenario with OAuth users in the database, you must 
-- either delete those users first or assign them a temporary hash before
-- running the ALTER below:
-- 
-- ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
