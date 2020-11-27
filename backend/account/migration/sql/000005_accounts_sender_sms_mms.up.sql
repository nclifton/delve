BEGIN;
ALTER TABLE
  account RENAME sender TO sender_sms;
ALTER TABLE
  account
ADD
  COLUMN sender_mms text [] DEFAULT '{""}';
DROP INDEX IF EXISTS account_sender;
CREATE INDEX account_sender_sms on account USING GIN(sender_sms);
CREATE INDEX account_sender_mms on account USING GIN(sender_mms);
COMMIT;
