-- name: CreateAccount :one
INSERT INTO accounts (
  owner, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateAccount :one
UPDATE accounts
  set owner = COALESCE(NULLIF(sqlc.arg(owner), ''), owner) ,
  balance = sqlc.arg(balance),
  currency = COALESCE(NULLIF(sqlc.arg(currency), ''), currency)
WHERE id = sqlc.arg(id)
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