CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE wallets (
                         id SERIAL PRIMARY KEY,
                         user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                         USD NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
                         RUB NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
                         EUR NUMERIC(18, 2) NOT NULL DEFAULT 0.00
);
