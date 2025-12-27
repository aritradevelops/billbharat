-- Create enum type "verification_type"
CREATE TYPE "public"."verification_type" AS ENUM ('email', 'phone', 'reset_password', 'change_email', 'change_phone');
-- Modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "password";
-- Create "passwords" table
CREATE TABLE "public"."passwords" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "created_by" uuid NOT NULL,
  "deleted_at" timestamptz NULL,
  "deleted_by" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "passwords_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "passwords_deleted_by_fkey" FOREIGN KEY ("deleted_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "passwords_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Modify "sessions" table
ALTER TABLE "public"."sessions" DROP CONSTRAINT "sessions_created_by_fkey", DROP CONSTRAINT "sessions_deleted_by_fkey", DROP CONSTRAINT "sessions_user_id_fkey", ADD COLUMN "expires_at" timestamptz NOT NULL, ADD CONSTRAINT "sessions_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "sessions_deleted_by_fkey" FOREIGN KEY ("deleted_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "sessions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Create "verification_requests" table
CREATE TABLE "public"."verification_requests" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "type" "public"."verification_type" NOT NULL,
  "code" character varying(16) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" uuid NOT NULL,
  "consumed_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "verification_requests_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "verification_requests_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
