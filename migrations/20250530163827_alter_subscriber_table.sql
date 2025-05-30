-- +goose Up
-- +goose StatementBegin
CREATE TABLE followers_new (
    id INTEGER PRIMARY KEY,
    display_name TEXT NOT NULL,
    username TEXT NOT NULL,
    followed_at INTEGER NOT NULL,
    followed_to TEXT NOT NULL REFERENCES users(id)
) STRICT;

INSERT INTO followers_new (id, display_name, username, followed_at, followed_to)
SELECT CAST(id AS INTEGER), display_name, username, followed_at, CAST(followed_to AS TEXT)
FROM followers;

DROP TABLE followers;

ALTER TABLE followers_new RENAME TO followers;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE followers_old (
    id TEXT NOT NULL PRIMARY KEY,
    display_name TEXT NOT NULL,
    username TEXT NOT NULL,
    followed_at INTEGER NOT NULL,
    followed_to INTEGER NOT NULL REFERENCES users(id)
) STRICT;


INSERT INTO followers_old (id, display_name, username, followed_at, followed_to)
SELECT CAST(id AS TEXT), display_name, username, followed_at, CAST(followed_to AS INTEGER)
FROM followers;

DROP TABLE followers;

ALTER TABLE followers_old RENAME TO followers;
-- +goose StatementEnd
