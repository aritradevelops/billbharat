CREATE TABLE "passwords" (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id uuid NOT NULL,
  password text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by uuid NOT NULL,
  deleted_at timestamptz,
  deleted_by uuid, 
  FOREIGN KEY (user_id) REFERENCES "users" (id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES "users" (id) ON DELETE CASCADE,
  FOREIGN KEY (deleted_by) REFERENCES "users" (id) ON DELETE CASCADE
);