-- +goose Up
-- +goose StatementBegin

ALTER TABLE events DROP COLUMN data;

CREATE TABLE conditions (
    id INTEGER PRIMARY KEY,
    broadcast_user_id TEXT,
    moderator_user_id TEXT,
    broadcaster_id TEXT,
    user_id TEXT
) STRICT;

ALTER TABLE events ADD COLUMN condition_id INTEGER REFERENCES conditions(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN condition_id;

DROP TABLE conditions;

ALTER TABLE events ADD COLUMN data TEXT NOT NULL;
-- +goose StatementEnd
