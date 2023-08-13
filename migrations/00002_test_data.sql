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
    tags (tag, author, create_time)
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

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM
    players
WHERE
    nickname = 'test';

-- +goose StatementEnd