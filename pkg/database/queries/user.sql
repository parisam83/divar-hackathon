-- name: InsertUser :execresult
INSERT INTO users (id)
VALUES ($1)
ON CONFLICT (id) DO NOTHING;