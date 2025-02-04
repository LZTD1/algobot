package domain

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"log"
	"os"
	"reflect"
	"testing"
	"tgbot/internal/domain"
	appError "tgbot/internal/error"
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
			truncateBase(base)

			user, err := sqlite3.User(1)
			if err == nil {
				fmt.Printf("%#v\n", user)
				t.Fatal("Expected error, got nil")
			}
		})
		t.Run("When user is created", func(t *testing.T) {
			truncateBase(base)

			base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(2, 'agent', 'cookie', 0);")
			base.Exec("INSERT INTO groups (group_id, owner_id, title, time_lesson) VALUES(0, ?, 'title', '2025-02-01 16:00:00');", 2)

			user, err := sqlite3.User(2)
			if err != nil {
				fmt.Printf("%s\n", err)
				t.Fatal("Expected user, got error")
			}

			assertUser(t, user, "agent", "cookie", false)
			assertGroups(
				t,
				user.Groups,
				[]domain.Group{
					{
						GroupID:    0,
						Title:      "title",
						TimeLesson: time.Date(2025, 2, 1, 16, 0, 0, 0, time.UTC),
					},
				},
			)
		})
	})
	t.Run("Test cookie method", func(t *testing.T) {
		t.Run("Cookie not exists", func(t *testing.T) {
			truncateBase(base)
			base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(3, 'agent', NULL, 0);")
			_, err := sqlite3.Cookie(3)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}
		})
		t.Run("Cookie exists", func(t *testing.T) {
			truncateBase(base)
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
		truncateBase(base)
		base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(5, 'agent', 'a', 0);")
		sqlite3.SetCookie(5, "cookie")
		cookie, _ := sqlite3.Cookie(5)

		if cookie != "cookie" {
			t.Errorf("Expected cookie to be 'cookie', got %s", cookie)
		}
	})
	t.Run("Test set userAgent", func(t *testing.T) {
		truncateBase(base)
		base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(6, 'a', 'a', 0);")
		sqlite3.SetUserAgent(6, "agent")

		row := base.QueryRow("SELECT u.user_agent FROM users u WHERE u.uid = ?", 6)

		var userAgent string
		row.Scan(&userAgent)

		if userAgent != "agent" {
			t.Errorf("Expected userAgent to be 'agent', got %s", userAgent)
		}
	})
	t.Run("Test Groups", func(t *testing.T) {
		t.Run("Get groups", func(t *testing.T) {
			truncateBase(base)
			base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(7, 'agent', 'cookie', 0);")
			base.Exec("INSERT INTO groups (group_id, owner_id, title, time_lesson) VALUES(0, ?, 'title', '2025-02-01 16:00:00');", 7)

			groups, err := sqlite3.Groups(7)
			if err != nil {
				t.Fatal(err)
			}

			assertGroups(
				t,
				groups,
				[]domain.Group{
					{
						GroupID: 0,
						Title:   "title",

						TimeLesson: time.Date(2025, 2, 1, 16, 0, 0, 0, time.UTC),
					},
				},
			)
		})
		t.Run("Set groups", func(t *testing.T) {
			truncateBase(base)
			wanted := []domain.Group{
				{
					GroupID: 0,
					Title:   "title",

					TimeLesson: time.Date(2025, 2, 1, 16, 0, 0, 0, time.UTC),
				},
			}

			base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(8, 'agent', 'cookie', 0);")
			sqlite3.SetGroups(8, wanted)

			groups, err := sqlite3.Groups(8)
			if err != nil {
				t.Fatal(err)
			}

			assertGroups(
				t,
				wanted,
				groups,
			)
		})
	})
	t.Run("Test notification method", func(t *testing.T) {
		t.Run("notification not exists", func(t *testing.T) {
			truncateBase(base)
			base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(3, 'agent', NULL, NULL);")
			notif, _ := sqlite3.Notification(3)
			if notif != false {
				t.Fatalf("Expected notification to be false, got true")
			}
		})
		t.Run("notification exists", func(t *testing.T) {
			truncateBase(base)
			base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(4, 'agent', 'cookie', 1);")
			notif, _ := sqlite3.Notification(4)
			if notif != true {
				t.Fatalf("Expected notification to be true, got false")
			}
		})
	})
	t.Run("Test set notification method", func(t *testing.T) {
		truncateBase(base)
		base.Exec("INSERT INTO users (uid, user_agent, cookie, notification) VALUES(5, 'agent', 'a', 0);")
		sqlite3.SetNotification(5, true)
		notif, _ := sqlite3.Notification(5)

		if notif != true {
			t.Errorf("Expected notif to be 'true', got %v", notif)
		}
	})
	t.Run("Register user ", func(t *testing.T) {
		truncateBase(base)

		sqlite3.RegisterUser(1)

		var exists bool
		base.QueryRow("SELECT EXISTS (SELECT 1 FROM users where uid = ?)", 1).Scan(&exists)
		if exists == false {
			t.Fatalf("User does not exist")
		}
	})
	t.Run("NotificationDate", func(t *testing.T) {
		t.Run("Get if not exists", func(t *testing.T) {
			truncateBase(base)
			base.Exec("INSERT INTO users (uid, user_agent, cookie, last_notification_msg, notification) VALUES(14, 'agent', 'cookie', '', 0);")
			_, err := sqlite3.LastNotificationDate(14)
			if err == nil {
				if errors.Is(err, appError.ErrNotValid) {
					t.Errorf("Wanted errNotValid, got %v", err)
				}
				t.Errorf("Wanted error, got nothing")
			}
		})
		t.Run("Get if exists", func(t *testing.T) {
			truncateBase(base)
			base.Exec("INSERT INTO users (uid, user_agent, cookie, last_notification_msg, notification) VALUES(14, 'agent', 'cookie', 'last_notification_msg', 0);")
			notif, err := sqlite3.LastNotificationDate(14)
			if err != nil {
				t.Errorf("Wanted result, got error")
			}
			if notif != "last_notification_msg" {
				t.Errorf("Wanted last_notification_msg, got %v", notif)
			}
		})
		t.Run("Set", func(t *testing.T) {
			truncateBase(base)
			base.Exec("INSERT INTO users (uid, user_agent, cookie, last_notification_msg, notification) VALUES(14, 'agent', 'cookie', 'last_notification_msg', 0);")
			err := sqlite3.SetLastNotificationDate(14, "asd")
			if err != nil {
				t.Errorf("Wanted change, got error")
			}
			notif, err := sqlite3.LastNotificationDate(14)
			if err != nil {
				t.Errorf("Wanted result, got error")
			}
			if notif != "asd" {
				t.Errorf("Wanted asd, got %v", notif)
			}
		})
	})
	t.Run("GetUsersByNotification", func(t *testing.T) {
		truncateBase(base)
		base.Exec("INSERT INTO users (uid, user_agent, cookie, last_notification_msg, notification) VALUES(14, 'agent', 'cookie', '', 0);")
		base.Exec("INSERT INTO users (uid, user_agent, cookie, last_notification_msg, notification) VALUES(15, 'agent', 'cookie', '', 1);")
		base.Exec("INSERT INTO users (uid, user_agent, cookie, last_notification_msg, notification) VALUES(16, 'agent', 'cookie', '', 0);")
		base.Exec("INSERT INTO users (uid, user_agent, cookie, last_notification_msg, notification) VALUES(17, 'agent', 'cookie', '', 1);")

		users, err := sqlite3.GetUsersByNotification(1)
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 2 {
			t.Errorf("Expected 2 users, got %v", len(users))
		}
		if users[0].UID != 15 || users[1].UID != 17 {
			t.Errorf("Expected 15, 17 UID user, got %v, %v", users[0].UID, users[0].UID)
		}
	})
}

func truncateBase(base *sql.DB) {
	base.Exec("DELETE FROM users;")
	base.Exec("DELETE FROM groups;")
}

func assertGroups(t *testing.T, groups []domain.Group, groups2 []domain.Group) {
	t.Helper()

	if !reflect.DeepEqual(groups, groups2) {
		t.Fatalf("Wanted equals, got %#v - %#v", groups, groups2)
	}
}

func assertUser(t *testing.T, user domain.User, userAgent string, cookie string, notif bool) {
	t.Helper()

	if user.UserAgent != userAgent {
		t.Fatalf("Expected userAgent to be %s, got %s", user.UserAgent, userAgent)
	}
	if user.Notifications != notif {
		t.Fatalf("Expected notifications to be %v, got %v", notif, user.Notifications)
	}
	if user.Cookie != cookie {
		t.Fatalf("Expected cookie to be %s, got %s", user.Cookie, cookie)
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
