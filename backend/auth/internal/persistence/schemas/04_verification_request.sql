-- migrations/001_verification_type_enum.sql
CREATE TYPE verification_type AS ENUM (
  'email',
  'phone',
  'reset_password',
  'change_email',
  'change_phone'
);

CREATE TABLE "verification_requests" (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id uuid NOT NULL,
  type verification_type NOT NULL,
  code VARCHAR(16) NOT NULL,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by uuid NOT NULL,
  consumed_at timestamptz,
  FOREIGN KEY (user_id) REFERENCES "users" (id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES "users" (id) ON DELETE CASCADE
);
  