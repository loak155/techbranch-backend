CREATE TABLE "bookmarks" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "article_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz,
  UNIQUE ("user_id", "article_id")
);

CREATE FUNCTION update_bookmark_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_bookmark
BEFORE UPDATE ON bookmarks
FOR EACH ROW
EXECUTE FUNCTION update_bookmark_updated_at();
