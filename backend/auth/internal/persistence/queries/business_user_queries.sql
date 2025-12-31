-- name: CreateBusinessUser :one
INSERT INTO "business_users" (
    user_id, business_id, role, created_by
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: FindBusinessesByUserID :many
SELECT sqlc.embed(bu), sqlc.embed(b) FROM "business_users" AS bu
LEFT JOIN "businesses" AS b ON bu.business_id = b.id AND b.deleted_at IS NULL
WHERE bu.user_id = $1 AND bu.deleted_at IS NULL;


-- name: FindUsersByBusinessID :many
SELECT sqlc.embed(bu), sqlc.embed(u) FROM "business_users" AS bu
LEFT JOIN "users" AS u ON bu.user_id = u.id AND u.deleted_at IS NULL
WHERE bu.business_id = $1 AND bu.deleted_at IS NULL;

-- name: DeleteBusinessUser :one
DELETE FROM "business_users" WHERE user_id = $1 AND business_id = $2 AND deleted_at IS NULL RETURNING *;