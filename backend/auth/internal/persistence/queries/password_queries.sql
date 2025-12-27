-- name: CreatePassword :exec
INSERT INTO "passwords" ("user_id", "password","created_by")
VALUES ($1, $2, $3);

-- name: DeletePassword :execrows
UPDATE "passwords" SET "deleted_at" = now(), "deleted_by" = $2 WHERE "user_id" = $1 AND "deleted_at" IS NULL;

-- name: FindPasswordByUserId :one
SELECT * FROM "passwords" WHERE "user_id" = $1 AND "deleted_at" IS NULL;

-- name: FindLastFourPasswordsByUserId :many
SELECT * FROM "passwords" WHERE "user_id" = $1 ORDER BY "created_at" DESC LIMIT 4;