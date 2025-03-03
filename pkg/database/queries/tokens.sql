-- name: AddToken :exec
INSERT INTO tokens 
(user_id, user_token, post_token, user_token_expiry, post_token_expiry)
VALUES ($1, $2, $3, $4, $5);
