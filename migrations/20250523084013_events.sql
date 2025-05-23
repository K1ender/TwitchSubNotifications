-- +goose Up
-- +goose StatementBegin
CREATE TABLE tokens (
    id INTEGER PRIMARY KEY,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at INTEGER NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
) STRICT;

CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    type TEXT NOT NULL,
    data TEXT NOT NULL
) STRICT;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE tokens;
DROP TABLE events;
-- +goose StatementEnd