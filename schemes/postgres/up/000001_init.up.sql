CREATE TABLE IF NOT EXISTS users
(
    id         uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    first_name varchar(255)     NOT NULL,
    last_name  varchar(255)     NOT NULL,
    email      varchar(255)     NOT NULL UNIQUE,
    wallet     uuid                      DEFAULT NULL,
    created_at timestamp        NOT NULL DEFAULT NOW(),
    updated_at timestamp        NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);

CREATE TABLE IF NOT EXISTS wallets
(
    id         uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    balance    numeric          NULL     DEFAULT NULL,
    reserved   numeric          NULL     DEFAULT NULL,
    created_at timestamp        NOT NULL DEFAULT NOW(),
    updated_at timestamp        NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions
(
    id         uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    wallet_id  uuid             NOT NULL,
    amount     numeric          NOT NULL,
    status     varchar(255)     NOT NULL,
    commentary text                      DEFAULT NULL,
    created_at timestamp        NOT NULL DEFAULT NOW(),
    updated_at timestamp        NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transaction_statuses
(
    id     serial PRIMARY KEY NOT NULL,
    status varchar(255)       NOT NULL UNIQUE
);

INSERT INTO transaction_statuses (status)
VALUES ('pending'),
       ('approved'),
       ('rejected');

ALTER TABLE users
    ADD CONSTRAINT users_wallet_fk FOREIGN KEY (wallet) REFERENCES wallets (id) ON DELETE CASCADE;

ALTER TABLE transactions
    ADD CONSTRAINT wallet_id_fk FOREIGN KEY (wallet_id) REFERENCES wallets (id) ON DELETE CASCADE,
    ADD CONSTRAINT status_fk FOREIGN KEY (status) REFERENCES transaction_statuses (status) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS services
(
    id            serial PRIMARY KEY NOT NULL,
    name          varchar(255)       NOT NULL,
    vendor_wallet uuid               NOT NULL,
    created_at    timestamp          NOT NULL DEFAULT NOW(),
    updated_at    timestamp          NOT NULL DEFAULT NOW()
);

ALTER TABLE services
    ADD CONSTRAINT vendor_wallet_fk FOREIGN KEY (vendor_wallet) REFERENCES wallets (id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS orders
(
    id         uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id    uuid             NOT NULL,
    service_id serial           NOT NULL,
    amount     numeric          NOT NULL,
    status     serial           NOT NULL,
    created_at timestamp        NOT NULL DEFAULT NOW(),
    updated_at timestamp        NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_statuses
(
    id     serial PRIMARY KEY NOT NULL,
    status varchar(255)       NOT NULL
);

INSERT INTO order_statuses (status)
VALUES ('created'),
       ('in_progress'),
       ('completed');

ALTER TABLE orders
    ADD CONSTRAINT status_fk FOREIGN KEY (status) REFERENCES order_statuses (id) ON DELETE CASCADE,
    ADD CONSTRAINT service_id_fk FOREIGN KEY (service_id) REFERENCES services (id) ON DELETE CASCADE,
    ADD CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

INSERT INTO wallets (id, balance, reserved)
VALUES ('1490b303-5e0e-476a-806c-7139a201c446', 0, 0),
       ('74231c9d-ede1-4e03-9b77-2133316b1771', 0, 0),
       ('aece5045-47ff-4ae5-99e6-9dd4bc4559cb', 0, 0);

INSERT INTO users (first_name, last_name, email, wallet)
VALUES ('Don', 'Don', 'dodo@ya.ru', '1490b303-5e0e-476a-806c-7139a201c446'),
       ('Apachai', 'Hopachai', 'kesha@rambler.ru', '74231c9d-ede1-4e03-9b77-2133316b1771'),
       ('Flipper', 'Zero', 'flipper@zero.ru', 'aece5045-47ff-4ae5-99e6-9dd4bc4559cb');

INSERT INTO services (name, vendor_wallet)
VALUES ('Dodo Pizza', '1490b303-5e0e-476a-806c-7139a201c446'),
       ('Yandex Taxi', '74231c9d-ede1-4e03-9b77-2133316b1771'),
       ('Yandex Food', 'aece5045-47ff-4ae5-99e6-9dd4bc4559cb');