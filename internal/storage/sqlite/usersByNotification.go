package sqlite

import (
	"algobot/internal/domain/models"
	"database/sql"
	"fmt"
)

func (s *Sqlite) UsersByNotification(wantNotif int) ([]models.User, error) {
	const op = "sqlite.UsersByNotification"

	q := `SELECT *  FROM users WHERE notification=?`
	pr, err := s.db.Prepare(q)
	if err != nil {
		return nil, fmt.Errorf("%s error while preparing sql: %w", op, err)
	}

	row, err := pr.Query(wantNotif)
	if err != nil {
		return nil, fmt.Errorf("%s error while Query sql: %w", op, err)
	}

	var users []models.User
	for row.Next() {
		var id int
		var uid int64
		var cookie sql.NullString
		var lastNotif sql.NullString
		var notif int

		if err := row.Scan(&id, &uid, &cookie, &lastNotif, &notif); err != nil {
			return nil, fmt.Errorf("%s error while scanning row: %w", op, err)
		}

		cookieNew := ""
		if cookie.Valid {
			cookieNew = cookie.String
		}
		lastNotifNew := ""
		if lastNotif.Valid {
			lastNotifNew = lastNotif.String
		}
		users = append(users, models.User{
			ID:               id,
			Uid:              uid,
			Cookie:           cookieNew,
			LastNotification: lastNotifNew,
			Notification:     notif,
		})

	}

	return users, nil
}
