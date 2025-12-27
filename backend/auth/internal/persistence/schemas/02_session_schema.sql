CREATE TABLE "sessions" (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  human_id VARCHAR(16) UNIQUE NOT NULL,
  user_ip VARCHAR(255) NOT NULL,
  user_agent VARCHAR(255) NOT NULL,
  refresh_token text NOT NULL UNIQUE,
  user_id uuid NOT NULL,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by uuid NOT NULL,
  deleted_at timestamptz,
  deleted_by uuid,
  FOREIGN KEY ("user_id") REFERENCES "users"("id"),
  FOREIGN KEY ("created_by") REFERENCES "users"("id"),
  FOREIGN KEY ("deleted_by") REFERENCES "users"("id")
);