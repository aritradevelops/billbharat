-- name: CreateUser :one
INSERT INTO "users" (
   human_id, name, email, email_verified,  phone, created_by
) VALUES (
   $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM "users" WHERE email = $1 AND deleted_at IS NULL;

-- name: FindUserById :one
SELECT * FROM "users" WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :one
UPDATE "users" SET deleted_at = CURRENT_TIMESTAMP AND deleted_by = $2 WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: DeactivateUser :one
UPDATE "users" SET deactivated_at = CURRENT_TIMESTAMP AND deactivated_by = $2 WHERE id = $1 AND deactivated_at IS NULL AND deleted_at IS NULL RETURNING *;

-- name: ActivateUser :one
UPDATE "users" SET deactivated_at = NULL AND updated_by = $2 WHERE id = $1 AND deactivated_at IS NOT NULL AND deleted_at IS NULL RETURNING *;

-- name: SetUserEmailVerified :one
UPDATE "users" SET email_verified = true WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: SetUserPhoneVerified :one
UPDATE "users" SET phone_verified = true WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: UpdateUserDP :one
UPDATE "users" SET dp = $2 WHERE id = $1 AND deleted_at IS NULL RETURNING *;
