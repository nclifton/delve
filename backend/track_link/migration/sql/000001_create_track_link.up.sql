BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS track_link (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  account_id UUID NOT NULL,
  message_id UUID NOT NULL,
  track_link_id TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  url text NOT NULL,
  hits INTEGER DEFAULT 0
);
CREATE INDEX track_link_id_idx ON track_link (account_id, track_link_id);
COMMIT;