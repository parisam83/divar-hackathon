-- name: AddPlace :exec
INSERT INTO places 
(name, type, latitude, longitude)
VALUES ($1, $2, $3, $4);
