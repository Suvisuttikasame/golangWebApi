CREATE TABLE "users" (
  "username" varchar NOT NULL PRIMARY KEY,
  "password" varchar NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY  ("owner") REFERENCES "users" ("username");
CREATE INDEX "owner_currency" ON "accounts" ("owner", "currency");
