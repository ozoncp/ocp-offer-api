-- +goose Up
-- +goose StatementBegin
CREATE TABLE "offer" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "team_id" BIGINT NOT NULL,
  "grade" BIGINT NOT NULL
);

CREATE UNIQUE INDEX "offer.user_team_id_index" ON "offer"("user_id", "team_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "offer";
-- +goose StatementEnd
