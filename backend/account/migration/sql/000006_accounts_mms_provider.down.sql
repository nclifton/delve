BEGIN;
ALTER TABLE
  account DROP COLUMN mms_provider_key;
COMMIT;
