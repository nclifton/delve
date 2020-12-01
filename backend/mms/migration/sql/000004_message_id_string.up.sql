BEGIN;
ALTER TABLE
  mms
ALTER COLUMN
  message_id
SET
  default '';
COMMIT;
