BEGIN;
ALTER TABLE
  account
ADD
  COLUMN sender text [] DEFAULT '{""}';
CREATE INDEX account_sender ON account (sender);
COMMIT;
