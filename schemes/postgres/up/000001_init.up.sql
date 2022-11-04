CREATE TABLE IF NOT EXISTS wallet_user
(
    id         uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    first_name varchar(255)     NOT NULL,
    last_name  varchar(255)     NOT NULL,
    email      varchar(255)     NOT NULL UNIQUE,
    wallet     uuid             NOT NULL,
    created_at timestamp        NOT NULL DEFAULT NOW(),
    updated_at timestamp        NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wallet
(
    id         uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    balance    decimal(10, 2)   NOT NULL DEFAULT 0,
    reserved   decimal(10, 2)   NOT NULL DEFAULT 0,
    created_at timestamp        NOT NULL DEFAULT NOW(),
    updated_at timestamp        NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS transaction_data
(
    id         uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    wallet_id  uuid             NOT NULL,
    amount     decimal(10, 2)   NOT NULL,
    status     serial           NOT NULL,
    commentary text             DEFAULT NULL,
    created_at timestamp        NOT NULL DEFAULT NOW(),
    updated_at timestamp        NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transaction_status
(
    id     serial PRIMARY KEY NOT NULL,
    status varchar(255)       NOT NULL
);

INSERT INTO transaction_status (status)
VALUES ('pending'),
       ('approved'),
       ('rejected');

ALTER TABLE wallet_user
    ADD CONSTRAINT wallet_user_wallet_fk FOREIGN KEY (wallet) REFERENCES wallet (id);

ALTER TABLE transaction_data
    ADD CONSTRAINT wallet_id_fk FOREIGN KEY (wallet_id) REFERENCES wallet (id),
    ADD CONSTRAINT status_fk FOREIGN KEY (status) REFERENCES transaction_status (id);