package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	appError "tgbot/internal/error"
	"time"
)

type Sqlite3 struct {
	db *sql.DB
}

func NewSqlite3(db *sql.DB) *Sqlite3 {
	return &Sqlite3{db: db}
}

func (s Sqlite3) GetUsersByNotification(notif int) ([]User, error) {
	row, err := s.db.Query(`SELECT u.uid, u.cookie, u.user_agent, u.notification FROM users u WHERE notification = ?`, notif)
	if err != nil {
		return nil, fmt.Errorf("NewSqlite3.GetUsersByNotification(%d) : %w", notif, row.Err())
	}

	var users []User
	for row.Next() {
		var baseId sql.NullInt64
		var cookie sql.NullString
		var userAgent sql.NullString
		var notifications bool

		err := row.Scan(&baseId, &cookie, &userAgent, &notifications)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("NewSqlite3.GetUsersByNotification(%d) : %w", notif, appError.ErrHasNone)
			}
			return nil, fmt.Errorf("NewSqlite3.GetUsersByNotification(%d) : %w", notif, row.Err())
		}

		u := User{
			UID:           baseId.Int64,
			Cookie:        cookie.String,
			UserAgent:     userAgent.String,
			Notifications: notifications,
			Groups:        nil,
		}
		users = append(users, u)
		s.appendGroups(&u, baseId.Int64)
	}

	return users, nil
}

func (s Sqlite3) LastNotificationDate(uid int64) (string, error) {
	row := s.db.QueryRow(`SELECT u.last_notification_msg FROM users u WHERE u.uid = ?`, uid)

	if row.Err() != nil {
		return "", fmt.Errorf("NewSqlite3.LastNotificationDate(%d) : %w", uid, row.Err())
	}

	var notifString sql.NullString

	err := row.Scan(&notifString)
	if err != nil {
		return "", fmt.Errorf("NewSqlite3.LastNotificationDate(%d) : %w", uid, err)
	}

	if !notifString.Valid || notifString.String == "" {
		return "", fmt.Errorf("NewSqlite3.LastNotificationDate(%d) : %w", uid, appError.ErrNotValid)
	}
	return notifString.String, nil
}

func (s Sqlite3) SetLastNotificationDate(uid int64, data string) error {
	_, err := s.db.Exec("UPDATE users SET last_notification_msg=? WHERE uid=?", data, uid)
	if err != nil {
		return fmt.Errorf("NewSqlite3.SetLastNotificationDate(%d, %s) : %w", uid, data, err)
	}
	return nil
}

func (s Sqlite3) User(uid int64) (User, error) {
	row := s.db.QueryRow(`SELECT u.uid, u.cookie, u.user_agent, u.notification FROM users u WHERE u.uid = ?`, uid)
	if row.Err() != nil {
		return User{}, fmt.Errorf("NewSqlite3.User(%d) : %w", uid, row.Err())
	}

	var baseId sql.NullInt64
	var cookie sql.NullString
	var userAgent sql.NullString
	var notifications bool

	err := row.Scan(&baseId, &cookie, &userAgent, &notifications)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("NewSqlite3.User(%d) : %w", uid, appError.ErrNotFound)
		}
		return User{}, fmt.Errorf("NewSqlite3.User(%d) : %w", uid, err)
	}

	if !baseId.Valid {
		return User{}, fmt.Errorf("NewSqlite3.User(%d) : %w", uid, appError.ErrNotValid)
	}

	u := User{
		Cookie:        cookie.String,
		UserAgent:     userAgent.String,
		Notifications: notifications,
	}

	s.appendGroups(&u, baseId.Int64)

	return u, nil
}

func (s Sqlite3) Cookie(uid int64) (string, error) {
	row := s.db.QueryRow(`SELECT u.cookie FROM users u WHERE u.uid = ?`, uid)

	if row.Err() != nil {
		return "", fmt.Errorf("NewSqlite3.Cookie(%d) : %w", uid, row.Err())
	}

	var cookie sql.NullString

	err := row.Scan(&cookie)
	if err != nil {
		return "", fmt.Errorf("NewSqlite3.Cookie(%d) : %w", uid, err)
	}

	if !cookie.Valid {
		return "", fmt.Errorf("NewSqlite3.Cookie(%d) : %w", uid, appError.ErrNotValid)
	}
	return cookie.String, nil
}

func (s Sqlite3) SetCookie(uid int64, cookie string) error {
	_, err := s.db.Exec("UPDATE users SET cookie=? WHERE uid=?", cookie, uid)
	if err != nil {
		return fmt.Errorf("NewSqlite3.SetCookie(%d, %s) : %w", uid, cookie, err)
	}
	return nil
}

func (s Sqlite3) SetUserAgent(uid int64, agent string) error {
	_, err := s.db.Exec("UPDATE users SET user_agent=? WHERE uid= ?;", agent, uid)
	if err != nil {
		return fmt.Errorf("NewSqlite3.SetUserAgent(%d, %s) : %w", uid, agent, err)
	}
	return nil
}

func (s Sqlite3) Groups(uid int64) ([]Group, error) {

	rows, err := s.db.Query(`SELECT g.group_id, g.title, g.time_lesson 
		FROM groups g 
		WHERE g.owner_id = ?;`, uid)
	if err != nil {
		return nil, fmt.Errorf("NewSqlite3.Groups(%d) : %w", uid, err)
	}
	defer rows.Close()

	groups := make([]Group, 0)
	for rows.Next() {
		var groupId sql.NullInt64
		var title sql.NullString
		var timeGroup sql.NullString

		if err := rows.Scan(&groupId, &title, &timeGroup); err != nil {
			return nil, fmt.Errorf("NewSqlite3.Groups(%d) : %w", uid, err)
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05", timeGroup.String)
		if err != nil {
			return nil, fmt.Errorf("NewSqlite3.Groups(%d) : %w", uid, err)
		}

		groups = append(groups, Group{
			GroupID:    int(groupId.Int64),
			Title:      title.String,
			TimeLesson: parsedTime,
		})
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("NewSqlite3.Groups(%d) : %w", uid, appError.ErrHasNone)
	}
	return groups, nil
}

func (s Sqlite3) SetGroups(uid int64, groups []Group) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("NewSqlite3.SetGroups(%d, %#v) : %w", uid, groups, err)
	}
	if _, err := tx.Exec("DELETE FROM groups WHERE owner_id=?", uid); err != nil {
		tx.Rollback()
		return fmt.Errorf("NewSqlite3.SetGroups(%d, %#v) : %w", uid, groups, err)
	}

	stmt, err := tx.Prepare("INSERT INTO groups (group_id, owner_id, title, time_lesson) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("NewSqlite3.SetGroups(%d, %#v) : %w", uid, groups, err)
	}
	defer stmt.Close()

	for _, g := range groups {
		if _, err := stmt.Exec(g.GroupID, uid, g.Title, g.TimeLesson.Format("2006-01-02 15:04:05")); err != nil {
			tx.Rollback()
			return fmt.Errorf("NewSqlite3.SetGroups(%d, %#v) : %w", uid, groups, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("NewSqlite3.SetGroups(%d, %#v) : %w", uid, groups, err)
	}

	return nil
}

func (s Sqlite3) Notification(uid int64) (bool, error) {
	row := s.db.QueryRow("SELECT notification FROM users WHERE uid=?", uid)

	var notif sql.NullBool
	if err := row.Scan(&notif); err != nil {
		return false, fmt.Errorf("NewSqlite3.Notification(%d) : %w", uid, err)
	}

	if !notif.Valid {
		return false, fmt.Errorf("NewSqlite3.Notification(%d) : %w", uid, appError.ErrNotValid)
	}

	return notif.Bool, nil
}
func (s Sqlite3) SetNotification(uid int64, notification bool) error {
	digit := 0
	if notification {
		digit = 1
	}

	_, err := s.db.Exec("UPDATE users SET notification=? WHERE uid=?", digit, uid)
	if err != nil {
		return fmt.Errorf("NewSqlite3.SetNotification(%d, %v) : %w", uid, notification, err)
	}
	return nil
}
func (s Sqlite3) RegisterUser(uid int64) error {
	_, err := s.db.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES (?, NULL, NULL, 0)", uid)
	if err != nil {
		return fmt.Errorf("NewSqlite3.RegisterUser(%d) : %w", uid, err)
	}
	return nil
}

func (s Sqlite3) appendGroups(u *User, id int64) {
	sqlQuery := "SELECT g.group_id, g.title, g.time_lesson FROM groups g WHERE g.owner_id = ?"
	query, err := s.db.Query(sqlQuery, id)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса, %s, со значением %v", sqlQuery, id)
		return
	}
	defer query.Close()
	for query.Next() {
		var groupId sql.NullInt64
		var title sql.NullString
		var timeGroup sql.NullString

		query.Scan(&groupId, &title, &timeGroup)

		parsedTime, err := time.Parse("2006-01-02 15:04:05", timeGroup.String)
		if err != nil {
			log.Printf("2 Ошибка при парсинге даты %v - %v", timeGroup.String, err)
			return
		}

		u.Groups = append(u.Groups, Group{
			GroupID:    int(groupId.Int64),
			Title:      title.String,
			TimeLesson: parsedTime,
		})
	}
}

func (s Sqlite3) Migrate(eFs fs.FS, dir string) {
	log.Println("Начинаю процесс миграции базы")
	files, err := fs.ReadDir(eFs, dir)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
        name TEXT PRIMARY KEY,
        executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		log.Fatalf("ошибка создания таблицы миграций: %v\n", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		log.Printf("Миграция файла - %s\n", file.Name())
		var exists bool
		err := s.db.QueryRow(
			`SELECT EXISTS (SELECT 1 FROM migrations WHERE name = ?)`,
			file.Name(),
		).Scan(&exists)
		if err != nil {
			log.Fatalf("Ошибка проверки таблицы: %v\n", err)
		}
		if exists {
			log.Printf("Миграция %s, уже существует\n", file.Name())
			continue
		}

		content, err := fs.ReadFile(eFs, "migrations/"+file.Name())
		if err != nil {
			log.Fatalf("ошибка чтения файла %s: %v", file.Name(), err)
		}
		_ = content

		tx, err := s.db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			log.Fatalf("ошибка выполнения миграции %s: %w", file.Name(), err)
		}

		if _, err := tx.Exec(
			"INSERT INTO migrations (name) VALUES (?)",
			file.Name(),
		); err != nil {
			tx.Rollback()
			log.Fatalf("ошибка записи миграции %s: %w", file.Name(), err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("ошибка коммита транзакции: %w", err)
		}

		log.Printf("Применена миграция: %s\n", file.Name())
	}

	log.Print("База данных готова к использованию!\n")
}
