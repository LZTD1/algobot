package domain

import (
	"database/sql"
	"errors"
	"io/fs"
	"log"
	"time"
)

type Sqlite3 struct {
	db *sql.DB
}

func NewSqlite3(db *sql.DB) *Sqlite3 {
	return &Sqlite3{db: db}
}

func (s Sqlite3) User(uid int64) (User, error) {
	sqlQuery := `SELECT u.id, u.cookie, u.user_agent, u.notification FROM users u WHERE u.uid = ?`
	row := s.db.QueryRow(sqlQuery, uid)
	if row.Err() != nil {
		log.Printf("Ошибка при выполнении запроса, %s, со значением %v", sqlQuery, uid)
		return User{}, row.Err()
	}

	var baseId sql.NullInt64
	var cookie sql.NullString
	var userAgent sql.NullString
	var notifications bool

	err := row.Scan(&baseId, &cookie, &userAgent, &notifications)
	if err != nil {
		return User{}, err
	}

	if !baseId.Valid {
		return User{}, errors.New("пользователь не найден")
	}

	u := User{
		cookie:        cookie.String,
		userAgent:     userAgent.String,
		notifications: notifications,
	}

	s.appendGroups(&u, baseId.Int64)

	return u, nil
}

func (s Sqlite3) Cookie(uid int64) (string, error) {
	sqlQuery := `SELECT u.cookie FROM users u WHERE u.uid = ?`
	row := s.db.QueryRow(sqlQuery, uid)

	if row.Err() != nil {
		return "", row.Err()
	}

	var cookie sql.NullString

	row.Scan(&cookie)

	if !cookie.Valid {
		return "", errors.New("cookie не установлены")
	}
	return cookie.String, nil
}

func (s Sqlite3) SetCookie(uid int64, cookie string) {
	_, err := s.db.Exec("UPDATE users SET cookie=? WHERE uid= ?;", cookie, uid)
	if err != nil {
		log.Printf("Ошибка при обновлении cookie [%v, %v] - %v", cookie, uid, err)
		return
	}
}

func (s Sqlite3) SetUserAgent(uid int64, agent string) {
	_, err := s.db.Exec("UPDATE users SET user_agent=? WHERE uid= ?;", agent, uid)
	if err != nil {
		log.Printf("Ошибка при обновлении useragent [%v, %v] - %v", agent, uid, err)
		return
	}
}

func (s Sqlite3) Groups(uid int64) ([]Group, error) {
	//TODO implement me
	panic("implement me")
}

func (s Sqlite3) SetGroups(uid int64, groups []Group) {
	//TODO implement me
	panic("implement me")
}

func (s Sqlite3) Notification(uid int64) bool {
	//TODO implement me
	panic("implement me")
}

func (s Sqlite3) RegisterUser(uid int64) {
	//TODO implement me
	panic("implement me")
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

		parsedTime, err := time.Parse("02.01.2006 15:04", timeGroup.String)
		if err != nil {
			log.Printf("Ошибка при парсинге даты %v - %v", timeGroup.String, err)
			return
		}

		u.groups = append(u.groups, Group{
			Id:   int(groupId.Int64),
			Name: title.String,
			Time: parsedTime,
		})
	}
}
