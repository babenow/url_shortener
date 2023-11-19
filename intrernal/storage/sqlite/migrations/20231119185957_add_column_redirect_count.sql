-- +goose Up
-- +goose StatementBegin
ALTER TABLE url
ADD COLUMN redirect_count INTEGER DEFAULT 0;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE url DROP COLUMN redirect_count;
-- +goose StatementEnd