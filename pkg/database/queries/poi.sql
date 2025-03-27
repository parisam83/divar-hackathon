-- name: GetOrCreatePoi :one
INSERT INTO poi(name, type, latitude, longitude)
VALUES ($1,$2,$3,$4)
ON CONFLICT (name, type) DO NOTHING
RETURNING *;

