BEGIN;
ALTER TABLE
  account
ADD
  COLUMN mms_provider_key text DEFAULT 'fake';
COMMIT;
