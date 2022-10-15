-- Sets now() time to updated_at column
CREATE
OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN
    NEW .updated_at = now();

RETURN NEW;

END;

$$language 'plpgsql';

-- Sets now() time to updated_at column on account update
CREATE TRIGGER update_account_update_time BEFORE
UPDATE
    ON account FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();