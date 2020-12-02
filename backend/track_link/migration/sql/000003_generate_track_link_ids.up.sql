BEGIN;
CREATE OR REPLACE FUNCTION track_link_id_generate(id_length INT)
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
        done := NOT EXISTS(SELECT 1 FROM track_link WHERE track_link_id = new_id);
    END LOOP;
    RETURN new_id;
END;
$$ LANGUAGE PLPGSQL
RETURNS NULL ON NULL INPUT
VOLATILE;
ALTER TABLE track_link ALTER COLUMN track_link_id SET DEFAULT track_link_id_generate(8);
COMMIT;
