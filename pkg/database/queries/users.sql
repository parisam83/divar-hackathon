-- name: GetUserPhoneNumberByID :one
SELECT phone_number FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserIDByPhoneNumber :one
SELECT id FROM users
WHERE phone_number = $1 LIMIT 1;

-- name: AddUser :exec
INSERT INTO users (phone_number)
VALUE ($1);
