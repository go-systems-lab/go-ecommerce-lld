-- Add email and password columns to accounts table
ALTER TABLE accounts ADD COLUMN email VARCHAR(255) UNIQUE;
ALTER TABLE accounts ADD COLUMN password VARCHAR(255);

-- Create index on email for faster login queries
CREATE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);

-- Add comments for documentation
COMMENT ON COLUMN accounts.email IS 'User email address for login';
COMMENT ON COLUMN accounts.password IS 'Hashed password for authentication';
