-- name: AddPOIResult :exec
INSERT INTO poi (post_token, place_id, distance, duration, snapp_cost, tapsi_cost)
VALUES ($1, $2, $3, $4, $5, $6);
