-- Remove email and password columns from accounts table
DROP INDEX IF EXISTS idx_accounts_email;
ALTER TABLE accounts DROP COLUMN IF EXISTS password;
ALTER TABLE accounts DROP COLUMN IF EXISTS email;
