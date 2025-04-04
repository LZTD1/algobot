package sqlite

import "fmt"

func (s *Sqlite) SetNotification(uid int64, isEnable bool) error {
	const op = "sqlite.SetNotification"

	digit := 0
	if isEnable {
		digit = 1
	}

	sqlq := "UPDATE users SET notification = ? WHERE uid = ?"
	pr, err := s.db.Prepare(sqlq)
	if err != nil {
		return fmt.Errorf("%s error while preparing statement: %w", op, err)
	}
	defer pr.Close()

	if _, err := pr.Exec(digit, uid); err != nil {
		return fmt.Errorf("%s error while executing statement: %w", op, err)
	}

	return nil
}
