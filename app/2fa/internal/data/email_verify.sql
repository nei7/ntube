-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (user_id, email, secret_code)
VALUES ($1, $2, $3)
RETURNING *;