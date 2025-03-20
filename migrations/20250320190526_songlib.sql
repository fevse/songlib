-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    band VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date TEXT,
    text TEXT,
    link TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd

