-- +goose Up
-- +goose StatementBegin

CREATE TABLE hubs (
  id BIGSERIAL PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW(),
  game_type TEXT NOT NULL,
  tg_id   bigint NOT NULL
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
