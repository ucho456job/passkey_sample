-- CreateTable
CREATE TABLE "public_key_credentials" (
    "credential_id" VARCHAR(255) NOT NULL,
    "user_id" UUID NOT NULL,
    "public_key" TEXT NOT NULL,
    "attestation_type" VARCHAR(50),
    "sign_count" INTEGER NOT NULL DEFAULT 0,
    "aagu_id" VARCHAR(36),
    "platform" TEXT NOT NULL,
    "user_agent" TEXT NOT NULL,
    "last_used_time" TIMESTAMP,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,

    CONSTRAINT "public_key_credentials_pkey" PRIMARY KEY ("credential_id")
);

-- CreateTable
CREATE TABLE "users" (
    "user_id" UUID NOT NULL,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,

    CONSTRAINT "users_pkey" PRIMARY KEY ("user_id")
);

-- CreateIndex
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

-- AddForeignKey
ALTER TABLE "public_key_credentials" ADD CONSTRAINT "public_key_credentials_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("user_id") ON DELETE CASCADE ON UPDATE CASCADE;
