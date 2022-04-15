INSERT INTO account (id, email, username, password, is_verified, is_archived)
VALUES ('d0fbc24f-5061-4d10-b92c-d386c8eba600', 'test@answersuck.com', 'test',
        '$2a$11$JgXOjzdX.1a.3lJciROxnuSkXFr43sMnWjLH59lctzMm84EtkOil.', TRUE, FALSE);

INSERT INTO account_verification_code (code, account_id)
VALUES ('5vkSjT9r6uSOLuwQzJ5xXlQ2pn5GBg5Zgcqd5dl9q6KDA2o2v9xfDBY74cgMq6KD',
        'd0fbc24f-5061-4d10-b92c-d386c8eba600');

INSERT INTO account_avatar (url, account_id)
VALUES ('https://avatars.dicebear.com/api/identicon/test.svg', 'd0fbc24f-5061-4d10-b92c-d386c8eba600');
