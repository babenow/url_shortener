-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url;
-- +goose StatementEnd