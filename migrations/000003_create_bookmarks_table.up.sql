CREATE TABLE "bookmarks" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "article_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

ALTER TABLE "bookmarks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "bookmarks" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id");
