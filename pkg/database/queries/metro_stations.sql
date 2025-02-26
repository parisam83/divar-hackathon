-- name: GetMetroStationName :one
SELECT name FROM metro_stations
WHERE id = $1 LIMIT 1;

-- name: AddMetroStation :exec
INSERT INTO metro_stations (name, latitude, longitude)
VALUES ($1, $2, $3);