-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (
    username, password, email
) VALUES (
    $1, $2, $3
)
RETURNING *;
