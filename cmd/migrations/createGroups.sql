CREATE TABLE groups
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id         INTEGER NOT NULL UNIQUE,
    owner_id         INTEGER,
    title            TEXT    NOT NULL,
    string_next_time TEXT    NOT NULL,
    time_lesson      TEXT    NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users (id)
);