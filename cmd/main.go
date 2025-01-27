package main

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"log"
	"tgbot/internal/domain"
)

const TOKEN = "6375608618:AAGtdaMkpj4SIJt495eNHOgw4oy5MZ_TIY4"

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	db, closeDb := getSqliteBase("base.db")
	defer closeDb()

	sqlite3 := domain.NewSqlite3(db)
	sqlite3.Migrate(migrationsFS, "migrations")
	//
	//pref := tele.Settings{
	//	Token:  TOKEN,
	//	Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	//}
	//
	//b, err := tele.NewBot(pref)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//defer b.Stop()
	//
	//boClient := clients.NewBackoffice("", clients.BackofficeSetting{
	//	Retry:        3,
	//	Timeout:      2 * time.Second,
	//	RetryTimeout: 1 * time.Second,
	//})
	//_ = boClient
	//svc := mocks.NewMockService(make(map[int64]bool))
	//state := stateMachine.NewMemory()
	//
	//regMid := middleware.NewRegister(svc)
	//
	//msgHandler := contextHandlers.NewOnText(svc, state)
	//callbackHandler := contextHandlers.NewOnCallback(svc, state)
	//
	//b.Use(regMid.Middleware, middleware.MessageLogger, middleware2.AutoRespond())
	//b.Handle(tele.OnText, msgHandler.Handle)
	//b.Handle(tele.OnCallback, callbackHandler.Handle)
	//
	//b.Start()
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
