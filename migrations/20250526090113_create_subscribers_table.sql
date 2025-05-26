-- +goose Up
-- +goose StatementBegin
CREATE TABLE followers (
    id TEXT NOT NULL PRIMARY KEY,
    display_name TEXT NOT NULL,
    username TEXT NOT NULL,
    followed_at INTEGER NOT NULL,
    followed_to INTEGER NOT NULL REFERENCES users(id)
) STRICT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE followers;
-- +goose StatementEnd
