BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS account (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS account_api_keys (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  account_id UUID NOT NULL,
  description text NOT NULL,
  key text NOT NULL,
  FOREIGN KEY (account_id) REFERENCES account(id)
);
CREATE INDEX account_api_keys_fkey ON account_api_keys (account_id);
COMMIT;
