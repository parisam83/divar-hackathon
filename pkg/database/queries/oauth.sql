-- name: AddOAuthData :exec
INSERT INTO oauth (session_id, access_token, refresh_token, expires_in, post_token)
VALUES ($1, $2, $3, $4, $5);

-- name: GetOAuthBySessionId :one
SELECT * FROM oauth WHERE session_id = $1;