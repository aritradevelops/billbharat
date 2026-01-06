-- name: CreateProductCategory :one
INSERT INTO "product_categories" ( name, business_id, created_by) 
VALUES ($1, $2, $3) RETURNING *;

-- name: SetProductCategoryNameByID :one
UPDATE "product_categories" 
SET name = $2, updated_at = now(), updated_by = $3
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: ListProductCategoriesByBusinessID :many
SELECT * FROM "product_categories" WHERE business_id = $1 AND deleted_at IS NULL ORDER BY name ASC LIMIT $2 OFFSET $3;

