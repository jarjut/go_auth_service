-- Create "users" table
CREATE TABLE "users" (
  "id" character varying(16) NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "name" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create "refresh_tokens" table
CREATE TABLE "refresh_tokens" (
  "id" bigserial NOT NULL,
  "user_id" character varying(16) NOT NULL,
  "token" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "is_revoked" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_refresh_tokens_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_refresh_tokens_token" to table: "refresh_tokens"
CREATE UNIQUE INDEX "idx_refresh_tokens_token" ON "refresh_tokens" ("token");
-- Create index "idx_refresh_tokens_user_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_user_id" ON "refresh_tokens" ("user_id");
