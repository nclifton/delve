BEGIN;

DROP TABLE IF EXISTS mms_task;
DROP TABLE IF EXISTS transcode_media;
DROP TYPE IF EXISTS media_type;
DROP TYPE IF EXISTS transcode_media_status;
DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;
