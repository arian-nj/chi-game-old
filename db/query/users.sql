-- name: GetTgUser :one
SELECT * FROM telegram_users
WHERE tg_id = $1;

-- name: GetAllTgUsers :many
SELECT * FROM telegram_users;

-- name: CountAllTgUsers :one
SELECT COUNT(*) FROM telegram_users;

-- name: CreateTgUser :one
INSERT INTO telegram_users (tg_id) VALUES ($1) ON
	CONFLICT (tg_id) DO UPDATE
	SET updated_at = NOW(), is_active = TRUE RETURNING *;

-- name: CountActiveTgUsers :one
SELECT COUNT(*) FROM telegram_users
WHERE is_active = TRUE;

-- name: CountUsersTgCreatedBetween :one
SELECT COUNT(*) FROM telegram_users
WHERE created_at >= $1 AND created_at <= $2;

--
-- -- name: ActiveTgUser :exec
-- UPDATE telegram_users SET is_active = TRUE
-- WHERE (
--   tg_id = $1
-- );
--
-- -- name: DiactiveTgUser :exec
-- UPDATE telegram_users SET is_active = FALSE
-- WHERE (
--   tg_id = $1
-- );
--

-- name: UpdateMixedTgUserStatuses :exec
UPDATE telegram_users
SET is_active = CASE
  WHEN id = ANY ($1::bigint[]) THEN TRUE
  ELSE FALSE  
END;


-- name: GetTgUsersStatic :one
SELECT
  COUNT(*) FILTER (WHERE is_active = TRUE) AS active_users,
  COUNT(*) AS total_users,
  COUNT(*) FILTER ( WHERE created_at >= NOW() - INTERVAL '1 day') AS users_created_last_24_hours,
  COUNT(*) FILTER ( WHERE created_at >= NOW() - INTERVAL '7 days') AS users_created_last_week,
  COUNT(*) FILTER ( WHERE created_at >= NOW() - INTERVAL '1 month') AS users_created_last_month
FROM telegram_users;
