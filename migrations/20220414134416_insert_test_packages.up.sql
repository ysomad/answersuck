INSERT INTO package (name, account_id, is_published, language_id)
VALUES ('Русский пакет 1', 'd0fbc24f-5061-4d10-b92c-d386c8eba600', FALSE, 1),
       ('Русский пакет 2', 'd0fbc24f-5061-4d10-b92c-d386c8eba600', FALSE, 1),
       ('Русский пакет 3', 'd0fbc24f-5061-4d10-b92c-d386c8eba600', TRUE, 1),
       ('English package 1', 'd0fbc24f-5061-4d10-b92c-d386c8eba600', FALSE, 2),
       ('English package 2', 'd0fbc24f-5061-4d10-b92c-d386c8eba600', FALSE, 2),
       ('English package 3', 'd0fbc24f-5061-4d10-b92c-d386c8eba600', TRUE, 2);

INSERT INTO package_tag (package_id, tag_id)
VALUES (1, 1),
       (1, 2),
       (1, 6),
       (1, 12),
       (2, 3),
       (2, 11),
       (2, 10),
       (3, 4),
       (3, 12),
       (3, 8),
       (3, 7),
       (4, 13),
       (4, 14),
       (4, 15),
       (5, 16),
       (5, 17),
       (5, 18),
       (5, 19),
       (6, 20),
       (6, 21),
       (6, 22),
       (6, 23),
       (6, 24);

INSERT INTO package_cover (url, package_id)
VALUES ('https://avatars.dicebear.com/api/identicon/rutest1.svg', 1),
       ('https://avatars.dicebear.com/api/identicon/rutest2.svg', 2),
       ('https://avatars.dicebear.com/api/identicon/rutest3.svg', 3),
       ('https://avatars.dicebear.com/api/identicon/entest1.svg', 4),
       ('https://avatars.dicebear.com/api/identicon/entest2.svg', 5),
       ('https://avatars.dicebear.com/api/identicon/entest3.svg', 6);

INSERT INTO stage (name, is_final, "order", package_id)
VALUES ('Этап 1', FALSE, 0, 1),
       ('Этап 2', FALSE, 1, 1),
       ('Этап 3', FALSE, 2, 1),
       ('Этап 4', FALSE, 3, 1),
       ('Финальный этап', TRUE, 5, 1),

       ('Этап 1', FALSE, 0, 2),
       ('Финал', TRUE, 5, 2),

       ('Этап 1', FALSE, 0, 3),
       ('Этап 2', FALSE, 1, 3),
       ('Этап 3', FALSE, 2, 3),
       ('Этап 4 финал', TRUE, 5, 3),

       ('Stage 1', FALSE, 0, 4),
       ('Stage 2', FALSE, 1, 4),
       ('Stage 3', FALSE, 2, 4),
       ('Stage 4', FALSE, 3, 4),
       ('Final stage!', TRUE, 5, 4),

       ('Stage 1', FALSE, 0, 5),
       ('Stage 2 final', TRUE, 5, 5),

       ('Stage 1', FALSE, 0, 6),
       ('Stage 2', FALSE, 1, 6),
       ('Final', TRUE, 5, 6);

INSERT INTO question_config (type, cost, interval, comment, secret_topic, secret_cost, is_keepable, is_visible)
VALUES ('DEFAULT', 100, 15, NULL, NULL, NULL, NULL, NULL),
       ('DEFAULT', 300, 15, NULL, NULL, NULL, NULL, NULL),
       ('DEFAULT', 400, 15, 'this question is sus', NULL, NULL, NULL, NULL),
       ('DEFAULT', 500, 15, NULL, NULL, NULL, NULL, NULL),
       ('DEFAULT', 1000, 30, 'host comment', NULL, NULL, NULL, NULL),

       ('BET', 150, 15, NULL, NULL, NULL, NULL, NULL),
       ('BET', 450, 15, NULL, NULL, NULL, NULL, NULL),
       ('BET', 1000, 30, 'bet comment', NULL, NULL, NULL, NULL),

       ('SECRET', 500, 15, NULL, 'secret topic 1', 5000, FALSE, FALSE),
       ('SECRET', 1000, 20, NULL, 'secret topic 2', 50, FALSE, FALSE),
       ('SECRET', 2000, 30, 'secret host comment', 'secret topic 3', 500, FALSE, FALSE),

       ('SUPERSECRET', 1000, 15, NULL, 'super secret topic 3', 10000, FALSE, TRUE),
       ('SUPERSECRET', 2000, 20, NULL, 'super secret topic 3', 500, TRUE, TRUE),
       ('SUPERSECRET', 3000, 30, 'super secret host comment', 'super secret topic 3', 50, TRUE, TRUE),
       ('SUPERSECRET', 5000, 45, NULL, 'super secret topic 3', 500, TRUE, FALSE),

       ('SAFE', 50, 15, NULL, NULL, NULL, NULL, NULL),
       ('SAFE', 150, 15, NULL, NULL, NULL, NULL, NULL),
       ('SAFE', 350, 20, NULL, NULL, NULL, NULL, NULL),
       ('SAFE', 500, 25, NULL, NULL, NULL, NULL, NULL),
       ('SAFE', 800, 40, NULL, NULL, NULL, NULL, NULL);

INSERT INTO stage_topic_question (stage_id, topic_id, question_id, question_config_id)
VALUES (1, 1, 1, 16),
       (1, 2, 2, 17),
       (1, 3, 3, 18),
       (2, 1, 4, 19),
       (2, 2, 5, 20),
       (2, 3, 6, 1),
       (3, 1, 7, 2),
       (3, 2, 8, 3),
       (3, 3, 9, 4),
       (4, 1, 10, 5),
       (4, 2, 11, 6),
       (4, 3, 12, 7),
       (5, 1, 1, 8),
       (5, 2, 2, 10),
       (5, 3, 3, 13),
       (5, 4, 4, 14),
       (5, 5, 5, 15),

       (6, 1, 1, 16),
       (6, 2, 2, 14),
       (6, 3, 3, 2),
       (6, 4, 4, 5),
       (6, 5, 5, 20),
       (7, 1, 6, 7),
       (7, 2, 7, 1),
       (7, 3, 8, 5),
       (7, 4, 9, 8),
       (7, 5, 10, 10),

       (8, 1, 1, 1),
       (8, 2, 2, 2),
       (8, 3, 3, 3),
       (9, 1, 4, 3),
       (9, 2, 5, 3),
       (9, 3, 6, 15),
       (10, 1, 7, 15),
       (10, 2, 8, 15),
       (10, 3, 9, 20),
       (11, 1, 10, 20),
       (11, 2, 11, 17),
       (11, 3, 12, 5),
       (11, 4, 1, 18),
       (11, 5, 2, 19),

       (12, 6, 13, 1),
       (12, 7, 14, 2),
       (12, 8, 15, 3),
       (13, 6, 16, 4),
       (13, 7, 17, 5),
       (13, 8, 18, 6),
       (14, 6, 19, 7),
       (14, 7, 20, 8),
       (14, 8, 21, 9),
       (15, 6, 22, 10),
       (15, 7, 23, 11),
       (15, 8, 24, 12),
       (16, 6, 13, 13),
       (16, 7, 14, 14),
       (16, 8, 15, 15),
       (16, 9, 16, 16),
       (16, 10, 17, 17),

       (17, 6, 13, 18),
       (17, 7, 14, 19),
       (17, 8, 15, 20),
       (17, 9, 16, 1),
       (17, 10, 17, 2),
       (18, 6, 18, 3),
       (18, 7, 19, 4),
       (18, 8, 20, 5),
       (18, 9, 21, 6),
       (18, 10, 22, 7),

       (19, 6, 13, 8),
       (19, 7, 14, 9),
       (19, 8, 15, 10),
       (19, 9, 16, 11),
       (20, 6, 17, 12),
       (20, 7, 18, 13),
       (20, 8, 19, 14),
       (20, 9, 20, 15),
       (21, 6, 21, 16),
       (21, 7, 22, 17),
       (21, 8, 23, 18);







