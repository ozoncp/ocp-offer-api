
-- Warnings:
-- Added the required column `is_deleted` to the `offer` table without a default value.
-- This is not possible if the table is not empty.

-- +goose Up
-- +goose StatementBegin
ALTER TABLE "offer" ADD COLUMN "is_deleted" BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE "offer" DROP COLUMN "is_deleted";
-- +goose StatementEnd
