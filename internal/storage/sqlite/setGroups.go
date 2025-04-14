package sqlite

import (
	"algobot/internal/domain/models"
	"fmt"
)

func (s *Sqlite) SetGroups(uid int64, groups []models.Group) error {
	const op = "sqlite.SetGroups"

	// Drop all
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s error while start tx: %w", op, err)
	}
	pr, err := tx.Prepare(`DELETE FROM groups WHERE owner_id = ?;`)
	if err != nil {
		return fmt.Errorf("%s error while preparing sql: %w", op, err)
	}
	_, err = pr.Exec(uid)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s error while exec sql: %w", op, err)
	}
	// Set new
	pr, err = tx.Prepare(`INSERT INTO groups (group_id, owner_id, title, time_lesson) VALUES (?, ?, ?, ?);`)
	if err != nil {
		return fmt.Errorf("%s error while preparing groups sql: %w", op, err)
	}
	for _, group := range groups {
		_, err = pr.Exec(group.GroupID, uid, group.Title, group.TimeLesson.Format("2006-01-02 15:04:05"))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s error while exec add group sql: %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s error while commit tx: %w", op, err)
	}

	return nil
}
