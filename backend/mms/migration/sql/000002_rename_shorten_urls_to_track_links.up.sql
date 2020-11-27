BEGIN;
ALTER TABLE mms RENAME COLUMN shorten_urls TO track_links;
COMMIT;