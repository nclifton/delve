BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS mms (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  account_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  provider_key TEXT NOT NULL,
  message_id TEXT,
  message_ref TEXT,
  country TEXT,
  subject TEXT,
  message TEXT,
  content_urls TEXT[],
  recipient TEXT NOT NULL,
  sender TEXT NOT NULL,
  status TEXT NOT NULL,
  shorten_urls BOOLEAN DEFAULT FAlSE,
  unsub BOOLEAN DEFAULT FALSE
);
COMMIT;