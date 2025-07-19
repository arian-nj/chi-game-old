-- +goose Up
-- +goose StatementBegin
ALTER TABLE users RENAME TO telegram_users
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
