-- name: CreatePassword :exec
INSERT INTO "passwords" ("user_id", "password","created_by")
VALUES ($1, $2, $3);

-- name: DeletePassword :execrows
DELETE FROM "passwords" WHERE "user_id" = $1 AND "deleted_at" IS NULL;

-- name: FindPasswordByUserId :one
SELECT * FROM "passwords" WHERE "user_id" = $1 AND "deleted_at" IS NULL;
