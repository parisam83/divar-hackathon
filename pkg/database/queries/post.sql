-- name: InsertPost :execresult
INSERT INTO posts (post_id, latitude, longitude,title)
VALUES ($1, $2, $3,$4)
ON CONFLICT (post_id) DO UPDATE
SET 
    latitude = EXCLUDED.latitude,
    longitude = EXCLUDED.longitude;
--     access_token = EXCLUDED.access_token,
--     refresh_token = EXCLUDED.refresh_token,
--     expires_in = EXCLUDED.expires_in,
--     updated_at=CURRENT_TIMESTAMP
-- WHERE now() > posts.expires_in;


-- name: UpdatePostCoordinates :execresult
UPDATE posts
SET 
    latitude = $1,
    longitude = $2
WHERE post_id = $3; 

-- name: GetPost :one
SELECT *
FROM posts
WHERE post_id = $1;


