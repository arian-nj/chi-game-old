-- name: CreateSessionMessage :one
INSERT INTO session_messages (session_id,player_id,message) VALUES ($1,$2,$3) RETURNING *;

-- name: GetSessionMessages :many
SELECT *
FROM session_messages
WHERE session_id = $1
ORDER BY id DESC
LIMIT $2;
