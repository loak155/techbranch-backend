CREATE TABLE "comments" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "article_id" bigint NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz
);

CREATE FUNCTION update_comment_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_comment
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_comment_updated_at();
