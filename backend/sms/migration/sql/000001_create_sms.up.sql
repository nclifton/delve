BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS sms (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  account_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  message_id TEXT,
  message_ref TEXT,
  country TEXT,
  message TEXT,
  sms_count SMALLINT,
  gsm BOOLEAN,
  recipient TEXT NOT NULL,
  sender TEXT NOT NULL,
  status TEXT NOT NULL
);
COMMIT;
