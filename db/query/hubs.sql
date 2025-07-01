-- name: CreateHub :one
INSERT INTO hubs (game_type, tg_id) VALUES ($1, $2) RETURNING *;
-- name: CountHubs :one
Select Count(*) from hubs;

-- name: CountLastHourHub :one
SELECT COUNT(*) AS hubs_last_hour
FROM hubs
WHERE created_at >= NOW() - INTERVAL '1 hour';

-- name: CountLastDayHubs :one
SELECT COUNT(*) AS hubs_today
FROM hubs
WHERE created_at >= date_trunc('day', NOW());
--
-- --name: CountLast7DaysHubs :many
-- WITH days AS (
--   SELECT generate_series(
--     date_trunc('day', NOW()) - interval '6 days',
--     date_trunc('day', NOW()),
--     interval '1 day'
--   ) AS day
-- )
-- SELECT
--   days.day,
--   COUNT(h.id) AS hubs_count
-- FROM days
-- LEFT JOIN hubs h
--   ON h.created_at >= days.day
--  AND h.created_at < days.day + interval '1 day'
-- GROUP BY days.day
-- ORDER BY days.day DESC;

