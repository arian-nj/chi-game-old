-- +goose Up
-- +goose StatementBegin
CREATE TABLE persons(
  id   BIGSERIAL PRIMARY KEY,
  tg_id   bigint NOT NULL UNIQUE,
  is_active boolean DEFAULT TRUE,
  updated_at timestamp NOT NULL DEFAULT NOW(),
  created_at timestamp NOT NULL DEFAULT NOW()
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

-- +goose StatementEnd
