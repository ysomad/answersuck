INSERT INTO answer_image (url)
VALUES ('https://avatars.dicebear.com/api/identicon/question_ru.svg'),
       ('https://avatars.dicebear.com/api/identicon/question_en.svg');

INSERT INTO answer (answer, answer_image_id)
VALUES ('Русский ответ 1', 1),
       ('Русский ответ 2', 1),
       ('Русский ответ 3', NULL),
       ('Русский ответ 4', NULL),
       ('Русский ответ 5', NULL),
       ('Русский ответ 6', 1),
       ('Русский ответ 7', 1),
       ('Русский ответ 8', 1),
       ('Русский ответ 9', NULL),
       ('Русский ответ 10', NULL),
       ('Русский ответ 11', NULL),
       ('Русский ответ 12', 1),
       ('English answer 1', 2),
       ('English answer 2', NULL),
       ('English answer 3', NULL),
       ('English answer 4', NULL),
       ('English answer 5', NULL),
       ('English answer 6', 2),
       ('English answer 7', 2),
       ('English answer 8', 2),
       ('English answer 9', 2),
       ('English answer 10', NULL),
       ('English answer 11', NULL),
       ('English answer 12', 2);

INSERT INTO question_media (url, type)
VALUES ('https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'VIDEO'),
       ('https://avatars.dicebear.com/api/identicon/question_media.svg', 'IMG'),
       ('https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3', 'AUDIO');

INSERT INTO question (question, answer_id, account_id, language_id, media_id)
VALUES ('Русский вопрос 1', 1, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 1),
       ('Русский вопрос 2', 2, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 2),
       ('Русский вопрос 3', 3, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 3),
       ('Русский вопрос 4', 4, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 2),
       ('Русский вопрос 5', 5, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 3),
       ('Русский вопрос 6', 6, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 1),
       ('Русский вопрос 7', 7, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 1),
       ('Русский вопрос 8', 8, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 1),
       ('Русский вопрос 9', 9, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 2),
       ('Русский вопрос 10', 10, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 3),
       ('Русский вопрос 11', 11, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 2),
       ('Русский вопрос 12', 12, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 1, 3),
       ('English question 1', 13, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 3),
       ('English question 2', 14, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 2),
       ('English question 3', 15, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 1),
       ('English question 4', 16, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 1),
       ('English question 5', 17, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 1),
       ('English question 6', 18, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 2),
       ('English question 7', 19, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 2),
       ('English question 8', 20, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 3),
       ('English question 9', 21, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 1),
       ('English question 10', 22, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 2),
       ('English question 11', 23, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 3),
       ('English question 12', 24, 'd0fbc24f-5061-4d10-b92c-d386c8eba600', 2, 1);

