-- name: InsertUser :exec
INSERT INTO users (id)
VALUES ($1)
ON CONFLICT (id) DO NOTHING;