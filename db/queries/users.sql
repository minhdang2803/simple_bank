-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users u 
WHERE u.username = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users AS u
SET 
    full_name = COALESCE($1,u.full_name),
    email = COALESCE($2,u.email)
WHERE u.username = $3;


-- name: ChangePassword :exec
UPDATE users AS u 
SET
    hashed_password = $1,
    password_changed_at = NOW()
WHERE u.username = $2;