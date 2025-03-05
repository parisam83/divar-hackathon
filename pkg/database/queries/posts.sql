-- name: AddPost :exec
INSERT INTO posts (post_token, latitude, longitude)
VALUES ($1, $2, $3);
