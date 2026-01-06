CREATE TABLE "users" (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  human_id VARCHAR(16) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  dp text,
  email_verified bool NOT NULL DEFAULT false,
  phone VARCHAR(16) NOT NULL,
  phone_verified bool NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by uuid NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_by uuid,
  deactivated_at timestamptz,
  deactivated_by uuid,
  deleted_at timestamptz,
  deleted_by uuid
  -- FOREIGN KEY (created_by) REFERENCES "users"("id"),
  -- FOREIGN KEY (updated_by) REFERENCES "users"("id"),
  -- FOREIGN KEY (deactivated_by) REFERENCES "users"("id"),
  -- FOREIGN KEY (deleted_by) REFERENCES "users"("id")
);