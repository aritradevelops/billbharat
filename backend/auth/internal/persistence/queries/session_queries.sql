-- name: CreateSession :exec
INSERT INTO "sessions" (
  human_id, user_id, user_ip, user_agent, refresh_token, expires_at, created_by
) VALUES (
   $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: FindSessionByRefreshToken :one
SELECT * FROM "sessions" WHERE refresh_token = $1 AND expires_at > CURRENT_TIMESTAMP AND deleted_at IS NULL;

-- name: DeleteSession :exec
UPDATE "sessions" SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;
