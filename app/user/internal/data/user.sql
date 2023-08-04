-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;
-- name: CreateUser :one
INSERT INTO users (email, password, username)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;