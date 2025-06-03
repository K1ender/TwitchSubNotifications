-- +goose Up
-- +goose StatementBegin
ALTER TABLE events RENAME TO events_old;

CREATE TABLE events (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    condition_id INTEGER REFERENCES conditions(id)
) STRICT;

INSERT INTO events (id, type, condition_id)
SELECT CAST(id AS TEXT), type, condition_id FROM events_old;

DROP TABLE events_old;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events RENAME TO events_old;

CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    type TEXT NOT NULL,
    condition_id INTEGER REFERENCES conditions(id)
) STRICT;

INSERT INTO events (id, type, condition_id)
SELECT CAST(id AS INTEGER), type, condition_id FROM events_old;

DROP TABLE events_old;
-- +goose StatementEnd
