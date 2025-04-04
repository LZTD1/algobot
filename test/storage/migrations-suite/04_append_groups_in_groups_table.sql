-- +goose Up
INSERT INTO groups (group_id, owner_id, title, time_lesson)
VALUES (1001, 999, 'group 1', '2025-03-23 14:00:00'),
       (1000, 999, 'group 2', '2025-03-23 16:00:00'),
       (999, 999, 'group 3', '2025-03-22 14:00:00');


-- +goose Down
DELETE
FROM groups
WHERE group_id in (1001, 1000, 999)