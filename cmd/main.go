package main

import (
	"algobot/internal_old/clients"
	"algobot/internal_old/contextHandlers"
	"algobot/internal_old/domain"
	"algobot/internal_old/middleware"
	"algobot/internal_old/schedulers"
	"algobot/internal_old/service"
	"algobot/internal_old/stateMachine"
	"database/sql"
	"embed"
	"github.com/joho/godotenv"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	tele "gopkg.in/telebot.v4"
	middleware2 "gopkg.in/telebot.v4/middleware"
	"log"
	"os"
	"time"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	// TODO зарефачить main
	godotenv.Load()

	TOKEN := os.Getenv("TELEGRAM_TOKEN")
	PORT := os.Getenv("GRPC_PORT")

	db, closeDb := getSqliteBase("base.db")
	defer closeDb()

	sqlite3 := domain.NewSqlite3(db)
	sqlite3.Migrate(migrationsFS, "migrations")

	pref := tele.Settings{
		Token:  TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer b.Stop()
	os.Setenv("TELEGRAM_NAME", b.Me.Username)

	boClient := clients.NewBackoffice("", clients.BackofficeSetting{
		Retry:        3,
		Timeout:      2 * time.Second,
		RetryTimeout: 1 * time.Second,
	})

	svc := service.NewDefaultService(sqlite3, boClient)
	state := stateMachine.NewMemory()

	regMid := middleware.NewRegister(svc)
	aiService := service.NewAiService(PORT)

	msgHandler := contextHandlers.NewOnText(svc, state, aiService)
	callbackHandler := contextHandlers.NewOnCallback(svc, state)

	tickerStop := goSchedule(b, svc)
	defer tickerStop()

	b.Use(regMid.Middleware, middleware.MessageLogger, middleware2.AutoRespond())
	b.Handle(tele.OnText, msgHandler.Handle)
	b.Handle(tele.OnCallback, callbackHandler.Handle)

	b.Start()
}

func goSchedule(b *tele.Bot, svc *service.DefaultService) func() {
	sch := schedulers.NewMessage(b, svc)
	ticker := time.NewTicker(10 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("SCHEDULER | Просмотр сообщений от детей ...")
				sch.Schedule()
			}
		}
	}()

	return ticker.Stop
}

func getSqliteBase(name string) (*sql.DB, func() error) {
	db, err := sql.Open("sqlite3", "file:"+name)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	log.Print("Подключение к базе данных установлено\n")
	return db, db.Close
}
