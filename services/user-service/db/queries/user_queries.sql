-- name: CreateUser :one
INSERT INTO users (
    username, email, password_hash, first_name, last_name
) VALUES (
    $1, $2, $3, $4, $5      --COALESCE($6, 'customer')
)
RETURNING id, username, email, first_name, last_name, is_active, created_at, updated_at;

--RETURNING *;

-- name: GetUserByID :one
SELECT id, username, email, password_hash, first_name, last_name, is_active, created_at, updated_at FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, first_name, last_name, is_active, created_at, updated_at FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT id, username, email, password_hash, first_name, last_name, is_active, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET username = $2,
    email = $3,
    first_name = $4,
    last_name = $5,
    is_active = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, username, email, password_hash, first_name, last_name, is_active, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
