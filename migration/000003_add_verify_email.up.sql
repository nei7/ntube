CREATE TABLE "verify_emails" (
    "id" bigserial PRIMARY KEY,
    "email" varchar NOT NULL,
    "user_id" UUID NOT NULL,
    "secret_code" varchar NOT NULL,
    "is_used" bool NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);
ALTER TABLE "verify_emails"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "users"
ADD COLUMN "is_email_verified" bool NOT NULL DEFAULT false;