BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION link_id_generate(id_length INT)
RETURNS TEXT AS $$
DECLARE
    new_id TEXT;
    done BOOL;
BEGIN
    done := false;
    WHILE NOT done LOOP
        new_id := array_to_string(
            ARRAY(
                SELECT substring(
                    'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789',
                    (random()*61)::int+1,
                    1
                )
                FROM generate_series(1,id_length) AS gs(x)
            )
            , ''
        );
        done := NOT EXISTS(SELECT 1 FROM opt_out WHERE link_id = new_id);
    END LOOP;
    RETURN new_id;
END;
$$ LANGUAGE PLPGSQL
RETURNS NULL ON NULL INPUT
VOLATILE LEAKPROOF;

CREATE TABLE IF NOT EXISTS opt_out (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  account_id UUID NOT NULL,
  message_id UUID NOT NULL,
  message_type TEXT NOT NULL,
  link_id TEXT DEFAULT link_id_generate(8),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

COMMIT;


