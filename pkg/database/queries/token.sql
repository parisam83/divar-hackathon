-- name: GetAccessTokenByUserIdPostId :one
SELECT * FROM tokens t
JOIN posts p ON p.post_id=t.post_id
JOIN users u ON u.id=t.user_id
WHERE  u.id = $1 AND 
        p.post_id = $2
LIMIT 1;