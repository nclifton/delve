BEGIN;

DROP TABLE IF EXISTS account;
DROP TABLE IF EXISTS account_api_keys;

DROP INDEX IF EXISTS conversation_log_fkey;

COMMIT;
