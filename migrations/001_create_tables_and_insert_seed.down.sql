-- DropForeignKey
ALTER TABLE "public_key_credentials" DROP CONSTRAINT "public_key_credentials_user_id_fkey";

-- DropIndex
DROP INDEX IF EXISTS "users_email_key";

-- DropTable
DROP TABLE IF EXISTS "public_key_credentials";

-- DropTable
DROP TABLE IF EXISTS "users";
