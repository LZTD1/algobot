-- +goose Up
INSERT INTO users (uid, cookie, last_notification_msg, notification)
VALUES (1001, 'cookie', NULL, 0),
       (1000, NULL, NULL, 0),
       (999, NULL, NULL, 0),
       (998, NULL, NULL, 1),
       (997, NULL, NULL, 2); -- 2 is only for test


-- +goose Down
DELETE
FROM users
WHERE uid in (1001, 1000, 999, 998, 997)