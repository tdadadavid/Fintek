-- name: CreateMoneyRecord :one
INSERT INTO money_records (
  user_id, 
  reference,
  amount,
  status
) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetMoneyRecordByReference :one
SELECT * FROM money_records WHERE reference = $1;

-- name: GetMoneyRecordByStatus :many
SELECT * FROM money_records WHERE status = $1;

-- name: DeleteMoneyRecordByID :exec
DELETE FROM money_records WHERE reference = $1;

