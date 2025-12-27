-- name: CreateVerificationRequest :exec
INSERT INTO "verification_requests" ("user_id", "type", "code", "expires_at", "created_by")
VALUES ($1, $2, $3, $4, $5);

-- name: FindVerificationRequestByUserIdAndType :one
SELECT * FROM "verification_requests" WHERE "user_id" = $1 AND "type" = $2 AND "expires_at" > CURRENT_TIMESTAMP AND "consumed_at" IS NULL;

-- name: SetConsumedAt :exec
UPDATE "verification_requests" SET "consumed_at" = CURRENT_TIMESTAMP WHERE "id" = $1 AND "consumed_at" IS NULL;