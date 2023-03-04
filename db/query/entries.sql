-- name: CreateEntry :one
INSERT INTO entries (
  account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: ListEntry :many
SELECT * FROM entries
ORDER BY id;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;