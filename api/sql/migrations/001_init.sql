-- -- +goose Up
-- CREATE OR REPLACE FUNCTION trigger_set_timestamp()
-- RETURNS trigger AS
-- '
-- BEGIN
--     NEW.updated_at = NOW();
-- 	RETURN NEW;
-- END;
-- '
-- LANGUAGE plpgsql;
--
-- -- +goose Down
-- DROP FUNCTION IF EXISTS trigger_set_timestamp;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
