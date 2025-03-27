-- name: InsertPost :exec
INSERT INTO posts (post_id, user_id, access_token, refresh_token, expires_in)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (post_id) DO UPDATE
SET 
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    expires_in = EXCLUDED.expires_in,
    updated_at=CURRENT_TIMESTAMP
WHERE now() > posts.expires_in;


-- name: UpdatePostCoordinates :exec
UPDATE posts
SET 
    latitude = $1,
    longitude = $2,
    coordinates_set = TRUE
WHERE post_id = $3; 