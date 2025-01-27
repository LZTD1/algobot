package domain

import (
	"database/sql"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"log"
	"os"
	"reflect"
	"testing"
	"tgbot/internal/domain"
	"time"
)

func TestDomain(t *testing.T) {
	const baseName = "temp.db"
	base, close := getSqliteBase(baseName)
	defer cleanup(baseName, close)

	sqlite3 := domain.NewSqlite3(base)
	sqlite3.Migrate(os.DirFS("../../cmd"), "migrations")

	t.Run("Test User method", func(t *testing.T) {
		t.Run("When user is not created", func(t *testing.T) {
			user, err := sqlite3.User(1)
			if err == nil {
				fmt.Printf("%#v\n", user)
				t.Fatal("Expected error, got nil")
			}
		})
		t.Run("When user is created", func(t *testing.T) {
			exec, _ := base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(2, 'agent', 'cookie', 0);")
			inseredId, _ := exec.LastInsertId()
			_, err := base.Exec("INSERT INTO groups (group_id, owner_id, title, string_next_time, time_lesson) VALUES(0, ?, 'title', 'string_next_time', '01.02.2025 16:00');", inseredId)
			if err != nil {
				fmt.Printf("%#v\n", err)
			}
			user, err := sqlite3.User(2)
			if err != nil {
				fmt.Printf("%s\n", err)
				t.Fatal("Expected user, got error")
			}

			assertUser(t, user, "agent", "cookie", false)
			assertGroups(
				t,
				user.Groups(),
				[]domain.Group{
					{
						Id:     0,
						Name:   "title",
						Lesson: "",
						Time:   time.Date(2025, 2, 1, 16, 0, 0, 0, time.UTC),
					},
				},
			)
		})
	})
	t.Run("Test cookie method", func(t *testing.T) {
		t.Run("Cookie not exists", func(t *testing.T) {
			base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(3, 'agent', NULL, 0);")
			_, err := sqlite3.Cookie(3)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}
		})
		t.Run("Cookie exists", func(t *testing.T) {
			base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(4, 'agent', 'cookie', 0);")
			cookie, err := sqlite3.Cookie(4)
			if err != nil {
				t.Error(err)
				t.Fatalf("Expected cookie, got error")
			}
			if cookie != "cookie" {
				t.Errorf("Expected cookie to be 'cookie', got %s", cookie)
			}
		})
	})
	t.Run("Test set cookie method", func(t *testing.T) {
		base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(5, 'agent', 'a', 0);")
		sqlite3.SetCookie(5, "cookie")
		cookie, _ := sqlite3.Cookie(5)

		if cookie != "cookie" {
			t.Errorf("Expected cookie to be 'cookie', got %s", cookie)
		}
	})
	t.Run("Test set userAgent", func(t *testing.T) {
		base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(6, 'a', 'a', 0);")
		sqlite3.SetUserAgent(6, "agent")

		row := base.QueryRow("SELECT u.user_agent FROM users u WHERE u.uid = ?", 6)

		var userAgent string
		row.Scan(&userAgent)

		if userAgent != "agent" {
			t.Errorf("Expected userAgent to be 'agent', got %s", userAgent)
		}
	})
}

func assertGroups(t *testing.T, groups []domain.Group, groups2 []domain.Group) {
	t.Helper()

	if !reflect.DeepEqual(groups, groups2) {
		t.Fatalf("Wanted equals, got %#v - %#v", groups, groups2)
	}
}

func assertUser(t *testing.T, user domain.User, userAgent string, cookie string, notif bool) {
	t.Helper()

	if user.UserAgent() != userAgent {
		t.Fatalf("Expected userAgent to be %s, got %s", user.UserAgent(), userAgent)
	}
	if user.Notifications() != notif {
		t.Fatalf("Expected notifications to be %v, got %v", notif, user.Notifications())
	}
	if user.Cookie() != cookie {
		t.Fatalf("Expected cookie to be %s, got %s", user.Cookie(), cookie)
	}
}

func cleanup(name string, closeBd func() error) {
	if err := closeBd(); err != nil {
		fmt.Printf("Ошибка при закрытии базы данных: %v\n", err)
	}

	if err := os.Remove(name); err != nil {
		fmt.Printf("Ошибка при удалении файла базы данных: %v\n", err)
	}
}

func getSqliteBase(name string) (*sql.DB, func() error) {
	db, err := sql.Open("sqlite3", "file:"+name)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	log.Print("Подключение к базе данных установлено\n")
	return db, db.Close
}
