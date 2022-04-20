CREATE TABLE IF NOT EXISTS stage_content
(
    id           serial        NOT NULL PRIMARY KEY,
    stage_id     int           NOT NULL REFERENCES stage (id),
    topic_id     int           NOT NULL REFERENCES topic (id),
    question_id  int           NOT NULL REFERENCES question (id),

    type         question_type NOT NULL DEFAULT 'DEFAULT',
    cost         int           NOT NULL,
    interval     int           NOT NULL,
    comment      text,

    secret_topic varchar(64),
    secret_cost  int,
    is_keepable  boolean,
    is_visible   boolean
);
