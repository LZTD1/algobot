package services

import (
	"reflect"
	"testing"
	"tgbot/internal/domain"
	appError "tgbot/internal/error"
	"tgbot/internal/models"
	"tgbot/internal/service"
	"tgbot/tests/mocks"
	"time"
)

func TestDefaultService(t *testing.T) {
	t.Run("Get cookie, with error", func(t *testing.T) {
		t.Run("Without error", func(t *testing.T) {
			d := mocks.MockDomain{}
			defaultService := service.NewDefaultService(&d, mocks.MockWebClient{})
			d.SetErrorCookie(nil)

			c, e := defaultService.Cookie(1)
			if e != nil {
				t.Fatalf("Wanted no error, got %v", e)
			}
			if c != "cookie" {
				t.Fatalf("Wanted 'cookie', got '%s'", c)
			}
		})
		t.Run("With error", func(t *testing.T) {
			d := mocks.MockDomain{}
			defaultService := service.NewDefaultService(&d, mocks.MockWebClient{})
			d.SetErrorCookie(appError.ErrNotValid)

			c, e := defaultService.Cookie(1)
			if e != nil {
				t.Fatalf("Wanted no error, got %v", e)
			}
			if c != "" {
				t.Fatalf("Wanted '', got %s", c)
			}
		})
	})
	t.Run("Get notification, with error", func(t *testing.T) {
		t.Run("Without error", func(t *testing.T) {
			d := mocks.MockDomain{}
			defaultService := service.NewDefaultService(&d, mocks.MockWebClient{})
			d.SetErrorNotif(nil)

			c, e := defaultService.Notification(1)
			if e != nil {
				t.Fatalf("Wanted no error, got %v", e)
			}
			if c != true {
				t.Fatalf("Wanted 'true', got '%v'", c)
			}
		})
		t.Run("With error", func(t *testing.T) {
			d := mocks.MockDomain{}
			defaultService := service.NewDefaultService(&d, mocks.MockWebClient{})
			d.SetErrorNotif(appError.ErrNotValid)

			c, e := defaultService.Notification(1)
			if e != nil {
				t.Fatalf("Wanted no error, got %v", e)
			}
			if c != false {
				t.Fatalf("Wanted 'false', got %v", c)
			}
		})
	})
	t.Run("Get CurrentGroup", func(t *testing.T) {
		defaultService := service.NewDefaultService(&mocks.MockDomain{}, mocks.MockWebClient{})
		group, err := defaultService.CurrentGroup(
			1,
			time.Date(2024, 10, 6, 14, 55, 55, 3, time.UTC),
		)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		wanted := models.Group{
			GroupID:    1,
			Title:      "test1",
			TimeLesson: time.Date(2024, 10, 6, 14, 55, 55, 3, time.UTC),
		}
		if !reflect.DeepEqual(wanted, group) {
			t.Fatalf("Wanted %v, got %v", wanted, group)
		}
	})
	t.Run("Get ActualInformation", func(t *testing.T) {
		domain := &mocks.MockDomain{}
		client := mocks.MockWebClient{}

		defaultService := service.NewDefaultService(domain, client)
		group, err := defaultService.ActualInformation(
			1,
			time.Date(2024, 10, 6, 14, 55, 55, 3, time.UTC),
			1,
		)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		wanted := models.ActualInformation{
			LessonTitle: "less3",
			LessonId:    3,
			MissingKids: []int{2},
		}
		if !reflect.DeepEqual(wanted, group) {
			t.Fatalf("Wanted %v, got %v", wanted, group)
		}
	})
	t.Run("Get AllKidsNames", func(t *testing.T) {
		domain := &mocks.MockDomain{}
		client := mocks.MockWebClient{}

		defaultService := service.NewDefaultService(domain, client)
		group, err := defaultService.AllKidsNames(
			1,
			1,
		)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		wanted := models.AllKids{
			1: models.KidData{
				FullName: "Иван Иванов",
				Login:    "ivan_ivanov",
				Password: "secret_password_123",
			},
			2: models.KidData{
				FullName: "Мария Петрова",
				Login:    "maria_petrov",
				Password: "password_321",
			},
		}
		if !reflect.DeepEqual(wanted, group) {
			t.Fatalf("Wanted %#v, got %#v", wanted, group)
		}
	})
	t.Run("Refresh groups", func(t *testing.T) {
		d := mocks.MockDomain{}
		defaultService := service.NewDefaultService(&d, mocks.MockWebClient{})
		err := defaultService.RefreshGroups(1)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}

		wanted := domain.Group{
			GroupID:    1,
			Title:      "Title",
			TimeLesson: time.Date(2025, 2, 1, 14, 00, 00, 00, time.UTC),
		}

		if !reflect.DeepEqual(d.MockGroups[0], wanted) {
			t.Fatalf("Wanted %#v, got %#v", wanted, d.MockGroups[0])
		}
	})
	t.Run("AllCredentials", func(t *testing.T) {
		d := mocks.MockDomain{}
		webClient := mocks.MockWebClient{}
		defaultService := service.NewDefaultService(&d, webClient)

		creds, err := defaultService.AllCredentials(1, 1)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}

		wanted := map[string]string{
			"Иван Иванов":   "ivan_ivanov:secret_password_123",
			"Мария Петрова": "maria_petrov:password_321",
		}

		if !reflect.DeepEqual(creds, wanted) {
			t.Fatalf("Wanted %#v, got %#v", wanted, creds)
		}
	})
	t.Run("UserUidsByNotif", func(t *testing.T) {
		d := mocks.MockDomain{}
		webClient := mocks.MockWebClient{}
		defaultService := service.NewDefaultService(&d, webClient)

		uids, err := defaultService.UsersByNotif(true)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}

		wanted := []models.ScheduleData{
			{
				UID:    1,
				Cookie: "2",
			},
		}
		if !reflect.DeepEqual(uids, wanted) {
			t.Fatalf("Wanted %#v, got %#v", wanted, uids)
		}
	})
	t.Run("NewMessageByUID", func(t *testing.T) {
		d := mocks.MockDomain{}
		webClient := mocks.MockWebClient{}
		defaultService := service.NewDefaultService(&d, webClient)

		msgs, err := defaultService.NewMessageByUID(1)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}

		wanted := []models.Message{
			{
				Id:      "2",
				From:    "2",
				Theme:   "2",
				Link:    "https://backoffice.algoritmika.org2",
				Content: "2",
			},
		}
		if !reflect.DeepEqual(msgs, wanted) {
			t.Fatalf("Wanted %#v, got %#v", wanted, msgs)
		}

		wantedDate := "29 дек. `24, 18:51"
		if d.DataNotif != wantedDate {
			t.Fatalf("Wanted %s, but got %s", wantedDate, d.DataNotif)
		}
	})
}
