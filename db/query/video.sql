-- name: CreateVideo :one
INSERT INTO videos (path, owner_id) VALUES ($1, $2) RETURNING *;

