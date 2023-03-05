-- name: CreateAccount :one
INSERT INTO accounts (
  owner, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateAccount :one
UPDATE accounts
  set owner = COALESCE(NULLIF($2, ''), owner),
  balance = $3,
  currency = COALESCE(NULLIF($4, ''), currency)
WHERE id = $1
RETURNING *;

-- name: ListAccount :many
SELECT * FROM accounts
ORDER BY id;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: DeleteAccount :one
DELETE FROM accounts
WHERE id = $1
RETURNING *;