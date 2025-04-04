package sqlite

import (
	"algobot/internal/domain/models"
	"database/sql"
	"fmt"
	"time"
)

func (s *Sqlite) Groups(uid int64) ([]models.Group, error) {
	const op = "sqlite.Groups"

	sqlq := "SELECT group_id, title, time_lesson FROM groups WHERE owner_id=?"
	pr, err := s.db.Prepare(sqlq)
	if err != nil {
		return nil, fmt.Errorf("%s error while preparing sql: %w", op, err)
	}
	defer pr.Close()

	rows, err := pr.Query(uid)
	if err != nil {
		return nil, fmt.Errorf("%s error while executing sql: %w", op, err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var groupId sql.NullInt64
		var title sql.NullString
		var timeGroup sql.NullString

		if err := rows.Scan(&groupId, &title, &timeGroup); err != nil {
			return nil, fmt.Errorf("%s error while scanning row: %w", op, err)
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05", timeGroup.String)
		if err != nil {
			return nil, fmt.Errorf("%s error while parsing time: %w", op, err)
		}

		groups = append(groups, models.Group{
			GroupID:    int(groupId.Int64),
			Title:      title.String,
			TimeLesson: parsedTime,
		})
	}

	return groups, nil
}
