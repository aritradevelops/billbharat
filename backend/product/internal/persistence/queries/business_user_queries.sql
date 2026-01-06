-- name: SyncBusinessUser :exec
INSERT INTO "business_users" (
    user_id,
    business_id,
    role,
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (user_id, business_id) DO 
UPDATE SET role = $3, created_at = $4, created_by = $5, updated_at = $6, updated_by = $7, deleted_at = $8, deleted_by = $9;
