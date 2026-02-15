CREATE TABLE app (
    id VARCHAR(36) NOT NULL PRIMARY KEY CHECK(length(id) > 0),
    content JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL CHECK(length(created_by) > 0),
    modified_by VARCHAR(50) NOT NULL CHECK(length(modified_by) > 0)
);

CREATE INDEX app_content_idx ON app USING GIN(content jsonb_path_ops);

CREATE OR REPLACE FUNCTION update_modified_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER app_update_modified_at
    BEFORE UPDATE ON app
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_at();