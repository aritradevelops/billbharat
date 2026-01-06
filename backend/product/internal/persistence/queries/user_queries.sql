-- name: SyncUser :exec
INSERT INTO users (
  id,
  human_id,
  name,
  email,
  dp,
  email_verified,
  phone,
  phone_verified,
  created_at,
  created_by,
  updated_at,
  updated_by,
  deactivated_at,
  deactivated_by,
  deleted_at,
  deleted_by
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) ON CONFLICT (id) DO 
UPDATE SET human_id = $2, name = $3, email = $4, dp = $5, email_verified = $6, phone = $7, 
 phone_verified = $8, created_at = $9,
 created_by = $10, updated_at = $11, updated_by = $12, deactivated_at = $13, deactivated_by = $14,
 deleted_at = $15, deleted_by = $16;
