CREATE TABLE IF NOT EXISTS stage_topic_question
(
    id                 serial NOT NULL PRIMARY KEY,
    stage_id           int    NOT NULL REFERENCES stage (id),
    topic_id           int    NOT NULL REFERENCES topic (id),
    question_id        int    NOT NULL REFERENCES question (id),
    question_config_id int    NOT NULL REFERENCES question_config (id)
);
