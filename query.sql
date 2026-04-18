-- name: GetAllGames :many
SELECT * FROM games;

-- name: InsertGame :one
INSERT INTO
    games (id, name, created_date)
VALUES (uuid_generate_v4 (), $1, $2) RETURNING *;