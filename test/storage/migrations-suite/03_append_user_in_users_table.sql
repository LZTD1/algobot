-- +goose Up
INSERT INTO users (uid, cookie, last_notification_msg, notification)
VALUES (1001, 'cookie', NULL, 0);


-- +goose Down
DELETE
FROM users
WHERE uid = 1001