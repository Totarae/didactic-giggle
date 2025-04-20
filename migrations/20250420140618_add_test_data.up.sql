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

---
UPDATE orders SET accrual = 100.0
WHERE order_number = '79927398713'; -- test1

UPDATE orders SET accrual = 300.5
WHERE order_number = '6011000990139424'; -- user_red

UPDATE orders SET accrual = 150.0
WHERE order_number = '5105105105105100'; -- user_blue

UPDATE orders SET accrual = 275.75
WHERE order_number = '3530111333300000';

-- test1 списал 40.5
INSERT INTO withdrawals (user_id, order_number, amount)
SELECT id, '79927398713', 40.5
FROM users WHERE login = 'test1';

-- user_blue списал 200
INSERT INTO withdrawals (user_id, order_number, amount)
SELECT id, '5105105105105100', 200.0
FROM users WHERE login = 'user_blue';