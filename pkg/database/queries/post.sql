-- -- name: InsertPost :execresult
-- INSERT INTO posts (post_id, latitude, longitude,title)
-- VALUES ($1, $2, $3,$4)
-- ON CONFLICT (post_id) DO UPDATE
-- SET 
--     latitude = EXCLUDED.latitude,
--     longitude = EXCLUDED.longitude;

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



-- name: InsertPost :execresult
INSERT INTO posts (
    post_id, 
    latitude, 
    longitude, 
    title, 
    owner_id
) 
VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT (post_id) 
DO UPDATE SET
    latitude = EXCLUDED.latitude,
    longitude = EXCLUDED.longitude,
    title = EXCLUDED.title,
    owner_id = COALESCE(posts.owner_id, EXCLUDED.owner_id)
RETURNING *;

-- name: CheckPostOwnership :one
SELECT 
    EXISTS (
        SELECT 1 
        FROM posts 
        WHERE owner_id = $1 AND post_id = $2
    ) AS isOwner;


-- name: UpdatePostOwner :execresult
UPDATE posts
SET owner_id = $1
WHERE post_id = $2 AND (owner_id IS NULL OR owner_id = $1);

-- name: GetPostDetails :one
SELECT * FROM posts
WHERE post_id = $1;