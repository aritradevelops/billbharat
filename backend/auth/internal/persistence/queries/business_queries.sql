-- name: CreateBusiness :one
INSERT INTO "businesses" (
    name, description, logo, industry, primary_currency, owner_id, currencies, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: FindBusinessById :one
SELECT * FROM "businesses" WHERE id = $1 AND deleted_at IS NULL;

-- name: FindBusinessByOwner :one
SELECT * FROM "businesses" WHERE owner_id = $1 AND deleted_at IS NULL;

-- name: UpdateBusiness :one
UPDATE "businesses" SET
    name = $2,
    description = $3,
    logo = $4,
    industry = $5,
    primary_currency = $6,
    currencies = $7,
    updated_by = $8
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SetBusinessLogo :one
UPDATE "businesses" SET
    logo = $2,
    updated_by = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;    

-- name: DeleteBusiness :one
DELETE FROM "businesses" WHERE id = $1 AND deleted_at IS NULL RETURNING *;