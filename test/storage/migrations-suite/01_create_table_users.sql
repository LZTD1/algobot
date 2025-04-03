-- +goose Up
CREATE TABLE users
(
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    uid                   INTEGER NOT NULL UNIQUE,
    cookie                TEXT    DEFAULT NULL,
    last_notification_msg TEXT    DEFAULT NULL,
    notification          INTEGER DEFAULT 0
);

-- +goose Down
DROP TABLE users;