-- name: CreateUser :one
INSERT INTO "Users" (
  email,
  hashed_password
) VALUES ($1, $2) RETURNING *;


-- name: GetUserByID :one
SELECT * FROM "Users" WHERE id=$1;

-- name: GetUserByEmail :one
SELECT * FROM "Users" WHERE email=$1;

-- name: ListUsers :many
SELECT * FROM "Users" ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUserPassword :one
UPDATE "Users" SET hashed_password = $1, updated_at=$2
WHERE id=$3 RETURNING *;

-- name: UpdateUserName :one
UPDATE "Users" SET username = $1, updated_at=$2
WHERE id=$3 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "Users" WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM "Users";