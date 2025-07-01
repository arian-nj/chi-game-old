-- name: GetUser :one
SELECT * FROM users
WHERE tg_id = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CountAllUsers :many
SELECT COUNT(*) FROM users;


-- name: CreateUser :exec
INSERT INTO users (tg_id) VALUES ($1)
ON CONFLICT (tg_id) DO UPDATE
SET updated_at = NOW(), is_active = TRUE;




    
-- name: ActiveUser :exec
UPDATE users SET is_active = TRUE
WHERE (
  tg_id = $1
);

-- name: DiactiveUser :exec
UPDATE users SET is_active = FALSE
WHERE (
  tg_id = $1
);

-- name: UpdateMixedUserStatuses :exec
UPDATE users
SET is_active = CASE
  WHEN id = ANY ($1::bigint[]) THEN TRUE
  ELSE FALSE  
END;


-- name: GetUsersStatic :one
SELECT
  COUNT(*) FILTER (WHERE is_active = TRUE) AS active_users,
  COUNT(*) AS total_users,
  COUNT(*) FILTER ( WHERE created_at >= NOW() - INTERVAL '1 day') AS users_created_last_24_hours,
  COUNT(*) FILTER ( WHERE created_at >= NOW() - INTERVAL '7 days') AS users_created_last_week,
  COUNT(*) FILTER ( WHERE created_at >= NOW() - INTERVAL '1 month') AS users_created_last_month
FROM users;
