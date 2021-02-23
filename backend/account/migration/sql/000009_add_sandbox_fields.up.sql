BEGIN;

CREATE TYPE account_mode AS ENUM (
  'live',
  'sandbox'
);

ALTER TABLE account ADD COLUMN mode account_mode NOT NULL DEFAULT 'live';
ALTER TABLE account ADD COLUMN parent_id UUID;


COMMIT;