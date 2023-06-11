CREATE TABLE users (
    user_id     uuid            PRIMARY KEY         DEFAULT gen_random_uuid(),
    username    VARCHAR(50)     NOT NULL,
    password    TEXT            NOT NULL,
    email       VARCHAR(255)    UNIQUE NOT NULL,
    enabled     BOOLEAN         NOT NULL            DEFAULT false,
    created_at  TIMESTAMP       NOT NULL            DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP
);
CREATE TABLE urls (
    url_id      uuid            PRIMARY KEY         DEFAULT gen_random_uuid(),
    url_key     VARCHAR(62)     UNIQUE NOT NULL,
    user_id     uuid            NOT NULL,
    created_at  TIMESTAMP       NOT NULL            DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);