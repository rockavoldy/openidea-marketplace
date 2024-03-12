DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'condition_enum') THEN
        CREATE TYPE "condition_enum" AS ENUM (
        'new',
        'second'
        );
    END IF;
END$$;


CREATE TABLE IF NOT EXISTS "products" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY UNIQUE PRIMARY KEY,
  "user_id" int NOT NULL,
  "name" varchar(255) NOT NULL,
  "price" decimal(10,2) NOT NULL DEFAULT 0,
  "image_url" varchar(255) NOT NULL,
  "stock" int NOT NULL DEFAULT 0,
  "condition" condition_enum NOT NULL DEFAULT 'new',
  "is_purchasable" bool DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

ALTER TABLE "products" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");