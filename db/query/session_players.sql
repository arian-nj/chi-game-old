-- name: CreateSessionPlayer :one
INSERT INTO session_players (session_id,person_id) VALUES ($1,$2) RETURNING *;
