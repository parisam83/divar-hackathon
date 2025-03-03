-- name: AddResult :exec
INSERT INTO results (token_id, metro_station_id, transportation_cost)
VALUES ($1, $2, $3);
