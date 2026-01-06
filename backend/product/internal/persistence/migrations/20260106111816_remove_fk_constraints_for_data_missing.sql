-- Modify "business_users" table
ALTER TABLE "public"."business_users" DROP CONSTRAINT "business_users_business_id_fkey", DROP CONSTRAINT "business_users_created_by_fkey", DROP CONSTRAINT "business_users_deleted_by_fkey", DROP CONSTRAINT "business_users_updated_by_fkey", DROP CONSTRAINT "business_users_user_id_fkey";
-- Modify "businesses" table
ALTER TABLE "public"."businesses" DROP CONSTRAINT "businesses_created_by_fkey", DROP CONSTRAINT "businesses_deleted_by_fkey", DROP CONSTRAINT "businesses_owner_id_fkey", DROP CONSTRAINT "businesses_updated_by_fkey";
-- Modify "product_categories" table
ALTER TABLE "public"."product_categories" DROP CONSTRAINT "product_categories_business_id_fkey", DROP CONSTRAINT "product_categories_created_by_fkey", DROP CONSTRAINT "product_categories_deleted_by_fkey", DROP CONSTRAINT "product_categories_updated_by_fkey";
-- Modify "users" table
ALTER TABLE "public"."users" DROP CONSTRAINT "users_created_by_fkey", DROP CONSTRAINT "users_deactivated_by_fkey", DROP CONSTRAINT "users_deleted_by_fkey", DROP CONSTRAINT "users_updated_by_fkey";
