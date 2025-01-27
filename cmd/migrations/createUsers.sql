CREATE TABLE users
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    uid          INTEGER NOT NULL UNIQUE,
    user_agent   TEXT    DEFAULT NULL,
    cookie       TEXT    DEFAULT NULL,
    notification INTEGER DEFAULT 0
);