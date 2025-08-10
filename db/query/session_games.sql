-- name: CreateSessionGame :one
INSERT INTO session_games (session_id,game_type) VALUES ($1,$2) RETURNING *;
