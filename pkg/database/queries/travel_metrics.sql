-- name: CreateTravelMetric :exec
INSERT INTO travel_metrics (distance, duration,origin_id, destionation_id)
VALUES ($1, $2, $3, $4);
