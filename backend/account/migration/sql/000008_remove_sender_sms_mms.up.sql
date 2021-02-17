BEGIN;

DROP INDEX IF EXISTS account_sender_sms;
DROP INDEX IF EXISTS account_sender_mms;

ALTER TABLE account DROP COLUMN mms_provider_key;
ALTER TABLE account DROP COLUMN sender_sms;
ALTER TABLE account DROP COLUMN sender_mms;

COMMIT;
