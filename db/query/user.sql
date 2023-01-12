-- name: GetUser :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    email,
    password
) VALUES (
    $1,
    $2
) RETURNING *;