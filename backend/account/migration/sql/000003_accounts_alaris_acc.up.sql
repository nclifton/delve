BEGIN;
ALTER TABLE
  account
ADD
  COLUMN alaris_username text default '',
ADD
  COLUMN alaris_password text default '',
ADD
  COLUMN alaris_url text default '';
COMMIT;
