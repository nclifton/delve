BEGIN;
CREATE TABLE IF NOT EXISTS webhook (
  id BIGSERIAL PRIMARY KEY,
  account_id TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  event TEXT NOT NULL,
  name TEXT NOT NULL,
  url TEXT NOT NULL,
  rate_limit INT NOT NULL
);
CREATE INDEX webhook_account_id on webhook(account_id);
COMMIT;
