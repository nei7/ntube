-- name: CreateVideo :one
INSERT INTO videos (path, owner_id, thumbnail, title, description) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetVideo :one
SELECT * FROM videos WHERE id = $1 LIMIT 1;

-- name: DeleteVideo :exec
DELETE FROM videos WHERE id = $1;

-- name: UpdateVideo :one
UPDATE videos SET 
  title = $1, description = $2 
  WHERE id = $3
  RETURNING *; 

-- name: GetUserVideos :many
SELECT * FROM videos WHERE owner_id = $1 ORDER BY uploaded_at LIMIT $2;

