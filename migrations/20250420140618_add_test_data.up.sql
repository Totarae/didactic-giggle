-- Включаем расширение pgcrypto (для crypt, gen_salt)
CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- Добавим пользователей
INSERT INTO users (login, password, password_hash)
VALUES
    ('test1', '123456', crypt('123456', gen_salt('bf'))),
    ('test2', '654321', crypt('654321', gen_salt('bf'))),
    ('user_blue', '123blue', crypt('123blue', gen_salt('bf'))),
    ('user_red', 'red345', crypt('red345', gen_salt('bf')));

-- Заказы для test1
INSERT INTO orders (order_number, user_id, status)
SELECT '79927398713', id, 'PROCESSED'
FROM users WHERE login = 'test1';

-- Заказы для test2
INSERT INTO orders (order_number, user_id, status)
SELECT '12345678903', id, 'PROCESSING'
FROM users WHERE login = 'test2';

-- Заказы для user_red (4 заказа)
INSERT INTO orders (order_number, user_id, status)
SELECT '4000000000000002', id, 'INVALID' FROM users WHERE login = 'user_red';
INSERT INTO orders (order_number, user_id, status)
SELECT '49927398716', id, 'NEW' FROM users WHERE login = 'user_red';
INSERT INTO orders (order_number, user_id, status)
SELECT '1234567812345670', id, 'PROCESSING' FROM users WHERE login = 'user_red';
INSERT INTO orders (order_number, user_id, status)
SELECT '6011000990139424', id, 'PROCESSED' FROM users WHERE login = 'user_red';

-- Заказы для user_blue (6 заказов)
INSERT INTO orders (order_number, user_id, status)
SELECT '4222222222222', id, 'NEW' FROM users WHERE login = 'user_blue';
INSERT INTO orders (order_number, user_id, status)
SELECT '4532015112830366', id, 'PROCESSING' FROM users WHERE login = 'user_blue';
INSERT INTO orders (order_number, user_id, status)
SELECT '5105105105105100', id, 'PROCESSED' FROM users WHERE login = 'user_blue';
INSERT INTO orders (order_number, user_id, status)
SELECT '378282246310005', id, 'INVALID' FROM users WHERE login = 'user_blue';
INSERT INTO orders (order_number, user_id, status)
SELECT '6011111111111117', id, 'NEW' FROM users WHERE login = 'user_blue';
INSERT INTO orders (order_number, user_id, status)
SELECT '3530111333300000', id, 'PROCESSED' FROM users WHERE login = 'user_blue';