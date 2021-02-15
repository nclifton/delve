BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS mms_task (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  provider_key TEXT NOT NULL,
  message_id TEXT,
  subject TEXT,
  message TEXT,
  sender TEXT NOT NULL,
  recipient TEXT NOT NULL,
  content_urls TEXT[],
  processed BOOLEAN NOT NULL DEFAULT false,
  failed BOOLEAN NOT NULL DEFAULT false
);

CREATE TYPE media_type AS ENUM (
  'IMAGE',
  'AUDIO',
  'VIDEO'
);

CREATE TYPE transcode_media_status AS ENUM (
  'READY',
  'PROCESSING',
  'COMPLETED',
  'FAILED'
);


CREATE TABLE IF NOT EXISTS transcode_media (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  media_type media_type,
  media_url TEXT NOT NULL,
  status transcode_media_status NOT NULL,
  external_transcode_id TEXT,
  external_status TEXT
);

COMMIT;