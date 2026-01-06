-- name: SyncBusiness :exec
INSERT INTO "businesses" (
    id,
    name,
    description,
    logo,
    industry,
    primary_currency,
    owner_id,
    currencies,
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) ON CONFLICT (id) DO 
UPDATE SET name = $2, description = $3, logo = $4, industry = $5, primary_currency = $6, owner_id = $7, currencies = $8, created_at = $9,
created_by = $10, updated_at = $11, updated_by = $12, deleted_at = $13, deleted_by = $14;