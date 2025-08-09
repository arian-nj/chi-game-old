CREATE TABLE session_games (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    game_id BIGINT NOT NULL REFERENCES games(id),
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP,
    status TEXT DEFAULT 'in_progress' -- could be finished, canceled, etc.
);
