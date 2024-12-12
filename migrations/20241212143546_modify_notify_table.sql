-- +goose Up
-- +goose StatementBegin
ALTER TABLE notify ADD create_at TIMESTAMP(0) DEFAULT NOW() NOT NULL;
ALTER TABLE notify ADD read_at TIMESTAMP(0) DEFAULT NULL::TIMESTAMP WITHOUT TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE notify DROP COLUMN create_at;
ALTER TABLE notify DROP COLUMN read_at;
-- +goose StatementEnd
