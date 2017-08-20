/* psql -d DATABASE < database.sql */

CREATE TABLE IF NOT EXISTS users
(
    id          serial NOT NULL,
    username    varchar(64) NOT NULL,
    created_at  timestamp NOT NULL DEFAULT now(),
    CONSTRAINT  pk_user_id primary key (id)
);

CREATE TABLE IF NOT EXISTS currencies
(
    id          serial NOT NULL,
    code        char(3) NOT NULL,
    rate        real NOT NULL,
    updated_at  timestamp NOT NULL DEFAULT now(),

    CONSTRAINT  pk_currency_id primary key (id),
    CONSTRAINT  currency_constraint UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS subscriptions
(
    id          serial NOT NULL,
    user_id     integer NOT NULL,
    currency_id integer NOT NULL,
    created_at  timestamp NOT NULL DEFAULT now(),

    CONSTRAINT  pk_subscription_id primary key (id),
    CONSTRAINT  user_currency_constraint UNIQUE (user_id, currency_id),
    CONSTRAINT  currency_fkey FOREIGN KEY (currency_id) REFERENCES currencies (id),
    CONSTRAINT  user_fkey FOREIGN KEY (user_id) REFERENCES users (id)
);