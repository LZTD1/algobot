package keyboards

import tele "gopkg.in/telebot.v4"

func Start() *tele.ReplyMarkup {
	StartKeyboard := &tele.ReplyMarkup{ResizeKeyboard: true}

	MissingBtn := StartKeyboard.Text("Получить отсутсвующих")
	MyGroupsBtn := StartKeyboard.Text("Мои группы")
	SettingsBtn := StartKeyboard.Text("Настройки")
	AIBtn := StartKeyboard.Text("AI 🔹")

	StartKeyboard.Reply(
		StartKeyboard.Row(MissingBtn),
		StartKeyboard.Row(MyGroupsBtn, SettingsBtn),
		StartKeyboard.Row(AIBtn),
	)

	return StartKeyboard
}
