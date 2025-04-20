CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    login         TEXT NOT NULL UNIQUE,
    password      TEXT NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE orders
(
    id           SERIAL PRIMARY KEY,
    order_number TEXT UNIQUE NOT NULL,
    user_id      INTEGER     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status       TEXT        NOT NULL DEFAULT 'NEW',
    accrual      NUMERIC,
    uploaded_at  TIMESTAMP DEFAULT NOW()
);