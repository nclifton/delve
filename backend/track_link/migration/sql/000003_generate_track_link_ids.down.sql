BEGIN;
ALTER TABLE track_link ALTER COLUMN track_link_id TYPE TEXT NOT NULL;
DROP FUNCTION IF EXISTS track_link_id_generate;
COMMIT;