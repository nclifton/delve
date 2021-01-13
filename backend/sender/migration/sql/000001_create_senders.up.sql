BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE provider_key AS ENUM (
  'fake',
  'optus',
  'mgage'
);
CREATE TABLE IF NOT EXISTS sender (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  account_id UUID,
  address TEXT NOT NULL,
  mms_provider_key TEXT,
  channels TEXT[],
  country TEXT NOT NULL,
  comment TEXT,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
CREATE INDEX sender_account_id on sender(account_id);

COMMIT;
