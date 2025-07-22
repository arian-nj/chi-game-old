-- name: GetPerson :one
SELECT * FROM persons
WHERE id = $1;

-- name: CreatePerson :one
INSERT INTO persons DEFAULT VALUES RETURNING *;

