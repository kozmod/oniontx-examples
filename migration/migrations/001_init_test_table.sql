-- +goose Up
CREATE TABLE IF NOT EXISTS text
(
    val text NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS text;
