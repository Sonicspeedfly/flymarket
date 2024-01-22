CREATE TABLE account (
    id BIGSERIAL PRIMARY KEY,
    phone TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    pass TEXT NOT NULL,
    roles TEXT[] NOT NULL DEFAULT '{}',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    image TEXT NOT NULL,
    category TEXT NOT NULL,
    file TEXT NOT NULL,
    information TEXT NOT NULL,
    count BIGINT NOT NULL,
    price BIGINT NOT NULL,
    account_id BIGINT NOT NULL REFERENCES account,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account_tokens (
    token TEXT NOT NULL UNIQUE,
    account_id BIGINT NOT NULL REFERENCES account,
    expire TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);