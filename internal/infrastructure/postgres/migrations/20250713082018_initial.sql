-- +goose Up
CREATE TABLE IF NOT EXISTS geckos
(
    id         VARCHAR(36) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    food_cycle VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS feeds
(
    date      VARCHAR(10)  NOT NULL,
    gecko_id  VARCHAR(36)  NOT NULL,
    food_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (date, gecko_id)
);

-- +goose Down
DROP TABLE IF EXISTS feeds;
DROP TABLE IF EXISTS geckos;
