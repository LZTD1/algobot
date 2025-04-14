-- +goose Up
CREATE TABLE groups
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id    INTEGER NOT NULL,
    owner_id    INTEGER,
    title       TEXT    NOT NULL,
    time_lesson TEXT    NOT NULL
);

-- +goose Down
DROP TABLE groups;