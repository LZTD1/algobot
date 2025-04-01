package keyboards

import tele "gopkg.in/telebot.v4"

func Start() *tele.ReplyMarkup {
	startKb := &tele.ReplyMarkup{ResizeKeyboard: true}

	missing := startKb.Text("Получить отсутсвующих")
	myGroups := startKb.Text("Мои группы")
	settings := startKb.Text("Настройки")
	ai := startKb.Text("AI 🔹")

	startKb.Reply(
		startKb.Row(missing),
		startKb.Row(myGroups, settings),
		startKb.Row(ai),
	)

	return startKb
}
