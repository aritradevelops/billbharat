-- Create "businesses" table
CREATE TABLE "public"."businesses" (
  "id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "logo" text NULL,
  "industry" character varying(50) NOT NULL,
  "primary_currency" character varying(10) NOT NULL,
  "owner_id" uuid NOT NULL,
  "currencies" text[] NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" uuid NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_by" uuid NULL,
  "deleted_at" timestamptz NULL,
  "deleted_by" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "businesses_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "businesses_deleted_by_fkey" FOREIGN KEY ("deleted_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "businesses_owner_id_fkey" FOREIGN KEY ("owner_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "businesses_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "business_users" table
CREATE TABLE "public"."business_users" (
  "user_id" uuid NOT NULL,
  "business_id" uuid NOT NULL,
  "role" character varying(50) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" uuid NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_by" uuid NULL,
  "deleted_at" timestamptz NULL,
  "deleted_by" uuid NULL,
  PRIMARY KEY ("user_id", "business_id"),
  CONSTRAINT "business_users_business_id_fkey" FOREIGN KEY ("business_id") REFERENCES "public"."businesses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "business_users_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "business_users_deleted_by_fkey" FOREIGN KEY ("deleted_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "business_users_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "business_users_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
