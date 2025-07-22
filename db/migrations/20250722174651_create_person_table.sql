-- +goose Up
-- +goose StatementBegin
CREATE TABLE persons (
  id   BIGSERIAL PRIMARY KEY,
	coins integer NOT NULL DEFAULT 0,
  created_at timestamp NOT NULL DEFAULT NOW()
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
