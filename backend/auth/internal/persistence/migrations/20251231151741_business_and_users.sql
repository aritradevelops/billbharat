-- Modify "businesses" table
ALTER TABLE "public"."businesses" ALTER COLUMN "id" SET DEFAULT gen_random_uuid();
