-- Modify "invitations" table
ALTER TABLE "public"."invitations" ADD COLUMN "role" character varying(255) NOT NULL;
