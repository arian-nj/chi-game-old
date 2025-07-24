-- +goose Up
-- +goose StatementBegin

CREATE TABLE game_sessions (
  id BIGSERIAL PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW()
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
