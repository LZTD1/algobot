package backoffice

import (
	"algobot/internal/config"
	"algobot/internal/domain/models"
	"algobot/internal/lib/backoffice"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestBackoffice(t *testing.T) {
	t.Run("Group", func(t *testing.T) {
		groupsExp := []models.Group{
			{GroupID: 98637162, Title: "Группа по курсу КГ", TimeLesson: time.Date(2025, time.April, 13, 14, 0, 0, 0, time.UTC)},
			{GroupID: 98623404, Title: "Группа по курсу ОЛиП МП", TimeLesson: time.Date(2025, time.April, 13, 12, 0, 0, 0, time.UTC)},
			{GroupID: 98621252, Title: "Группа по курсу Пст", TimeLesson: time.Date(2025, time.April, 13, 18, 0, 0, 0, time.UTC)},
			{GroupID: 98619913, Title: "Группа по курсу КГ", TimeLesson: time.Date(2025, time.April, 12, 10, 0, 0, 0, time.UTC)},
			{GroupID: 98619873, Title: "Группа по курсу ГД", TimeLesson: time.Date(2025, time.April, 12, 14, 0, 0, 0, time.UTC)},
			{GroupID: 98619867, Title: "Группа по курсу ВП", TimeLesson: time.Date(2025, time.April, 12, 12, 0, 0, 0, time.UTC)},
			{GroupID: 98589447, Title: "Группа по курсу ВП", TimeLesson: time.Date(2025, time.April, 13, 10, 0, 0, 0, time.UTC)},
			{GroupID: 985504, Title: "Группа по курсу Пст 2", TimeLesson: time.Date(2025, time.April, 12, 18, 0, 0, 0, time.UTC)},
			{GroupID: 978298, Title: "Группа по курсу Пст 2", TimeLesson: time.Date(2025, time.April, 13, 16, 0, 0, 0, time.UTC)},
		}

		responseHTML := readFile("group_example")

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(responseHTML))
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		bo := backoffice.NewBackoffice(&config.Backoffice{
			Retries:         5,
			RetriesTimeout:  time.Second,
			ResponseTimeout: time.Second,
		}, backoffice.WithURL(server.URL))

		group, err := bo.Group("cookie")
		assert.NoError(t, err)
		assert.Equal(t, groupsExp, group)
	})
}

func readFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	responseHTML := string(b)
	return responseHTML
}
