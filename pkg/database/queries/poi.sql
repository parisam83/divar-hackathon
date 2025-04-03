-- name: GetOrCreatePoi :one
INSERT INTO poi(name, type, latitude, longitude)
VALUES ($1,$2,$3,$4)
ON CONFLICT (name, type) DO NOTHING
RETURNING *;

-- name: UpsertSubwayStation :one
INSERT INTO poi (name, type, latitude, longitude)
VALUES ($1, 'subway', $2, $3)
ON CONFLICT (name, type) 
DO UPDATE SET 
    latitude = EXCLUDED.latitude,
    longitude = EXCLUDED.longitude
RETURNING id;

-- name: GetToSubwayInfo :one
-- Get the nearest subway station to a given lat,lng (There is an issue, many posts can have the same lat,lng)
SELECT 
    p.name AS station_name,
    tm.distance,
    tm.duration
FROM travel_metrics tm
JOIN poi p ON p.id = tm.destination_id
JOIN posts pt ON pt.post_id = tm.origin_id
WHERE 
    pt.latitude = $1 AND 
    pt.longitude = $2 AND
    p.type = 'subway'
LIMIT 1;

