BEGIN;
DROP INDEX IF EXISTS account_api_keys_key;
CREATE INDEX account_api_keys_key on account_api_keys(key);
DROP INDEX IF EXISTS account_sender;
CREATE INDEX account_sender on account USING GIN(sender);
COMMIT;
