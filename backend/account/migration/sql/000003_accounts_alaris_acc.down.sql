BEGIN;
ALTER TABLE
  account DROP COLUMN alaris_username,
  DROP COLUMN alaris_password,
  DROP COLUMN alaris_url;
COMMIT;
