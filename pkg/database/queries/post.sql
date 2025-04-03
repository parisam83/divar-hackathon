-- name: InsertPost :execresult
INSERT INTO posts (post_id, latitude, longitude,title)
VALUES ($1, $2, $3,$4)
ON CONFLICT (post_id) DO NOTHING;
-- SET 
--     access_token = EXCLUDED.access_token,
--     refresh_token = EXCLUDED.refresh_token,
--     expires_in = EXCLUDED.expires_in,
--     latitude = EXCLUDED.latitude,
--     longitude = EXCLUDED.longitude,
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


-- name: InsertToken :execresult
INSERT INTO tokens (post_id, user_id, access_token, refresh_token, expires_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (post_id, user_id) DO UPDATE
SET
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    expires_at= EXCLUDED.expires_at
WHERE now() > tokens.expires_at;