-- +goose Up
-- +goose StatementBegin
CREATE TABLE session_messages (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    player_id BIGINT NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    sent_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
