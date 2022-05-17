CREATE TABLE IF NOT EXISTS package_tag
(
    package_id int NOT NULL REFERENCES package (id),
    tag_id     int NOT NULL REFERENCES tag (id),
    PRIMARY KEY (package_id, tag_id)
);
