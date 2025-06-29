-- +migrate Up
CREATE TABLE IF NOT EXISTS pets
(
    id         VARCHAR(36) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    food_cycle VARCHAR(255) DEFAULT ''
);

CREATE TABLE IF NOT EXISTS feedings
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    date      VARCHAR(10) NOT NULL UNIQUE,
    pet_id    VARCHAR(36) NOT NULL,
    food_type VARCHAR(255)
);

-- +migrate Down
DROP TABLE IF EXISTS feedings;
DROP TABLE IF EXISTS pets;
