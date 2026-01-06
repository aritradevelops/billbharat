-- name: CreateInvitation :one
INSERT INTO "invitations" (
  name, email, phone, role, business_id, hash, expires_at, created_by
) VAlUES (
  $1,$2,$3,$4,$5,$6,$7,$8
) RETURNING *;

-- name: FindInvitationByHash :one
SELECT * FROM "invitations" WHERE hash = $1 AND expires_at > now() AND deleted_at IS NULL;

-- name: SetInvitationAcceptedAt :exec
UPDATE "invitations" SET accepted_at = CURRENT_TIMESTAMP WHERE id = $1 AND accepted_at IS NULL AND deleted_at IS NULL;
