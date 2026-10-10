CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 1000
    account_type TEXT NOT NULL DEFAULT 'user'
        CHECK (account_type IN ('user', 'admin', 'seller'))
);