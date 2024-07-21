CREATE TABLE "articles" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "url" text NOT NULL,
  "image" text,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);