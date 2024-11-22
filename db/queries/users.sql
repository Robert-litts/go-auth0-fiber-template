-- name: CreateUser :one
INSERT INTO users (auth0_id, email)
VALUES ($1, $2)
RETURNING id, auth0_id, email, created_at;

-- name: GetUserByAuth0ID :one
SELECT id, auth0_id, email, created_at
FROM users
WHERE auth0_id = $1;
