-- +goose Up
-- +goose StatementBegin
INSERT INTO
    players(
        nickname,
        email,
        display_name,
        email_verified,
        password,
        create_time
    )
VALUES
    (
        'test',
        'test@test.com',
        'test player',
        true,
        '$argon2id$v=19$m=65536,t=1,p=2$wvbSeXYgL6cAzha1qGu16w$aMD4W28gyZ51CI52WgbWxWCiHfYuRwv8nXL8m2eJ7CA',
        '2023-08-03 22:25:08.947 +0700'
    );

INSERT INTO
    media (url, type, uploader, create_time)
VALUES
    (
        'https://upload.wikimedia.org/wikipedia/en/0/0d/Van_Halen_album.jpg',
        1,
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        'https://upload.wikimedia.org/wikipedia/en/9/9a/Trollface_non-free.png',
        1,
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        'https://media.tenor.com/eUJ9ASaUTXQAAAAC/boomer-meme.gif',
        1,
        'test',
        '2023-08-03 22:25:08.947 +0700'
    );

INSERT INTO
    packs (
        id,
        name,
        author,
        is_published,
        cover_url,
        create_time
    )
VALUES
    (
        1337,
        'test pack',
        'test',
        false,
        'https://upload.wikimedia.org/wikipedia/en/0/0d/Van_Halen_album.jpg',
        '2023-08-03 22:25:08.947 +0700'
    );

INSERT INTO
    tags (name, author, create_time)
VALUES
    (
        'видеоигры',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        'мемы',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        'фильмы',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    );

INSERT INTO
    pack_tags(pack_id, tag)
VALUES
    (1337, 'видеоигры'),
    (1337, 'мемы'),
    (1337, 'фильмы');

INSERT INTO
    rounds(id, name, position, pack_id)
VALUES
    (1, 'Раунд 1', 1, 1337),
    (2, 'Раунд 2', 2, 1337),
    (3, 'Финал', 3, 1337);

INSERT INTO
    topics (id, title, author, create_time)
VALUES
    (
        10,
        'хорроры',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        20,
        'рпг',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        30,
        'ходилки',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        40,
        'шутеры',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        50,
        'аркада',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        60,
        'история мемов',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        70,
        'видео мемы',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        80,
        'мешапы',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        90,
        'RYTP',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        100,
        'фэнтези',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        110,
        'ужасы',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        120,
        'артахус',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    ),
    (
        130,
        'трейлеры',
        'test',
        '2023-08-03 22:25:08.947 +0700'
    );

INSERT INTO
    round_topics(id, round_id, topic_id)
VALUES
    (1, 1, 10),
    (2, 1, 20),
    (3, 1, 30),
    (4, 1, 40),
    (5, 1, 50),
    (6, 2, 60),
    (7, 2, 70),
    (8, 2, 80),
    (9, 2, 90),
    (10, 3, 100),
    (11, 3, 110),
    (12, 3, 120),
    (13, 3, 130);

INSERT INTO
    answers(id, text, media_url)
VALUES
    (
        1337,
        'Ответ на все вопросы!',
        'https://upload.wikimedia.org/wikipedia/en/9/9a/Trollface_non-free.png'
    );

INSERT INTO
    questions(
        id,
        text,
        answer_id,
        author,
        media_url,
        create_time
    )
VALUES
    (
        228,
        'Тестовый вопрос',
        1337,
        'test',
        'https://media.tenor.com/eUJ9ASaUTXQAAAAC/boomer-meme.gif',
        '2023-08-03 22:25:08.947 +0700'
    );

INSERT INTO
    round_questions(
        id,
        round_topic_id,
        question_id,
        question_type,
        cost,
        answer_time,
        host_comment,
        secret_topic,
        secret_cost,
        transfer_type,
        is_keepable
    )
VALUES
    (
        111,
        1,
        228,
        1,
        100,
        15000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        222,
        1,
        228,
        2,
        300,
        10000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        333,
        2,
        228,
        2,
        50,
        5000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        444,
        2,
        228,
        3,
        300,
        15000000000,
        'хост ишак',
        'Секретная тема вопроса!',
        100,
        NULL,
        FALSE
    ),
    (
        555,
        3,
        228,
        3,
        500,
        15000000000,
        'хост ишак',
        'СУПЕР Секретная тема вопроса!',
        1000,
        2,
        TRUE
    ),
    (
        666,
        3,
        228,
        1,
        100,
        15000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        777,
        4,
        228,
        1,
        100,
        15000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        888,
        4,
        228,
        2,
        300,
        10000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        999,
        5,
        228,
        2,
        50,
        5000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        1111,
        5,
        228,
        3,
        300,
        15000000000,
        'хост ишак',
        'Секретная тема вопроса!',
        100,
        NULL,
        FALSE
    ),
    (
        1222,
        6,
        228,
        3,
        500,
        15000000000,
        'хост ишак',
        'СУПЕР Секретная тема вопроса!',
        1000,
        2,
        TRUE
    ),
    (
        1333,
        6,
        228,
        1,
        100,
        15000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        1444,
        7,
        228,
        1,
        100,
        15000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        1555,
        7,
        228,
        2,
        300,
        10000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        1666,
        8,
        228,
        2,
        50,
        5000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        1777,
        8,
        228,
        3,
        300,
        15000000000,
        'хост ишак',
        'Секретная тема вопроса!',
        100,
        NULL,
        FALSE
    ),
    (
        1888,
        9,
        228,
        3,
        500,
        15000000000,
        'хост ишак',
        'СУПЕР Секретная тема вопроса!',
        1000,
        2,
        TRUE
    ),
    (
        1999,
        9,
        228,
        1,
        100,
        15000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        2111,
        10,
        228,
        1,
        100,
        15000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        2222,
        10,
        228,
        2,
        300,
        10000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        3333,
        11,
        228,
        2,
        50,
        5000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        4444,
        11,
        228,
        3,
        300,
        15000000000,
        'хост ишак',
        'Секретная тема вопроса!',
        100,
        NULL,
        FALSE
    ),
    (
        5555,
        12,
        228,
        3,
        500,
        15000000000,
        'хост ишак',
        'СУПЕР Секретная тема вопроса!',
        1000,
        2,
        TRUE
    ),
    (
        6666,
        12,
        228,
        1,
        100,
        15000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        7777,
        13,
        228,
        3,
        500,
        15000000000,
        'хост ишак',
        'СУПЕР Секретная тема вопроса!',
        1000,
        2,
        TRUE
    ),
    (
        8888,
        13,
        228,
        1,
        100,
        15000000000,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM
    players
WHERE
    nickname = 'test';

-- +goose StatementEnd