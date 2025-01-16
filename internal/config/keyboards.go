package config

import tele "gopkg.in/telebot.v4"

var (
	StartKeyboard    = &tele.ReplyMarkup{ResizeKeyboard: true}
	SettingsKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}

	missing  = StartKeyboard.Text("Получить отсутсвующих")
	myGroups = StartKeyboard.Text("Мои группы")
	settings = StartKeyboard.Text("Настройки")
)

func init() {
	StartKeyboard.Reply(
		StartKeyboard.Row(missing),
		StartKeyboard.Row(myGroups, settings),
	)
}
