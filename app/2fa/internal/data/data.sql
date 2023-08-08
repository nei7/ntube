-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (user_id, email, secret_code)
VALUES ($1, $2, $3)
RETURNING *;
-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used = TRUE
WHERE id = @id
    AND secret_code = @secret_code
    AND is_used = FALSE
    AND expired_at > now()
RETURNING *;
-- name: UpdateUser :one
UPDATE users
SET is_email_verified = sqlc.arg(is_email_verified)
WHERE id = sqlc.arg(id)
RETURNING *;