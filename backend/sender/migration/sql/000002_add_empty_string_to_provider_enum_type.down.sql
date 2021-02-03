BEGIN;

ALTER TYPE provider_key RENAME TO _provider_key;
CREATE TYPE provider_key AS ENUM (
  'fake',
  'optus',
  'mgage'
);
ALTER TABLE sender RENAME COLUMN mms_provider_key TO _mms_provider_key;
ALTER TABLE sender ADD mms_provider_key provider_key;
UPDATE sender SET mms_provider_key = _mms_provider_key::text::provider_key;
ALTER TABLE sender DROP COLUMN _mms_provider_key;
DROP TYPE _mms_provider_key;

COMMIT;
