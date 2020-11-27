BEGIN;
ALTER TABLE
  account RENAME sender_sms TO sender;
ALTER TABLE
  account DROP COLUMN sender_mms;
DROP INDEX IF EXISTS account_sender_sms;
DROP INDEX IF EXISTS account_sender_mms;
CREATE INDEX account_sender on account USING GIN(sender);
COMMIT;
