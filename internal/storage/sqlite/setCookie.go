package sqlite

import "fmt"

func (s *Sqlite) SetCookie(uid int64, cookie string) error {
	const op = "sqlite.SetCookie"

	sqlq := "UPDATE users SET cookie=? WHERE uid=?"
	pr, err := s.db.Prepare(sqlq)
	if err != nil {
		return fmt.Errorf("%s error while preparing statement: %w", op, err)
	}
	defer pr.Close()

	if _, err := pr.Exec(cookie, uid); err != nil {
		return fmt.Errorf("%s error while executing statement: %w", op, err)
	}

	return nil
}
