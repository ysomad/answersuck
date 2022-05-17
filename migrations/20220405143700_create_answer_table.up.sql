CREATE TABLE IF NOT EXISTS answer_image
(
    id  serial        NOT NULL PRIMARY KEY,
    url varchar(2048) NOT NULL
);

CREATE TABLE IF NOT EXISTS answer
(
    id              serial       NOT NULL PRIMARY KEY,
    answer          varchar(100) NOT NULL,
    answer_image_id int REFERENCES answer_image (id)
);
