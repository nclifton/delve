BEGIN;
ALTER TABLE mms RENAME COLUMN track_links TO shorten_urls;
COMMIT;