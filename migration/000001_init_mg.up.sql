CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE wallets (
                         user_id INT PRIMARY KEY,
                         usd NUMERIC DEFAULT 0.0,
                         rub NUMERIC DEFAULT 0.0,
                         eur NUMERIC DEFAULT 0.0,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              user_id INT NOT NULL,
                              amount NUMERIC NOT NULL,
                              currency VARCHAR(10) NOT NULL,
                              type VARCHAR(20) NOT NULL,
                              timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              status VARCHAR(20) NOT NULL,
                              FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE exchange_rates (
                                from_currency VARCHAR(10) NOT NULL,
                                to_currency VARCHAR(10) NOT NULL,
                                rate NUMERIC NOT NULL,
                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                PRIMARY KEY (from_currency, to_currency)
);
