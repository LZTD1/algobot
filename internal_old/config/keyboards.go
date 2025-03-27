package config

import tele "gopkg.in/telebot.v4"

var (
	StartKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}
	MissingBtn    = StartKeyboard.Text("Получить отсутсвующих")
	MyGroupsBtn   = StartKeyboard.Text("Мои группы")
	SettingsBtn   = StartKeyboard.Text("Настройки")
	AIBtn         = StartKeyboard.Text("AI 🔹")

	MyGroupsKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}
	refreshGroupsBtn = MyGroupsKeyboard.Data("Обновить группы", "refresh_groups")

	SettingsKeyboard      = &tele.ReplyMarkup{ResizeKeyboard: true}
	SetCookieBtn          = SettingsKeyboard.Data("Установить Cookie", "set_cookie")
	ChangeNotificationBtn = SettingsKeyboard.Data("Переключить уведомления", "change_notification")

	RejectKeyboard  = &tele.ReplyMarkup{ResizeKeyboard: true}
	RejectActionBtn = RejectKeyboard.Text("Отменить действие")

	AIKeyboard      = &tele.ReplyMarkup{ResizeKeyboard: true}
	BackBtn         = StartKeyboard.Text("⬅️ Назад")
	ClearHistoryBtn = StartKeyboard.Text("Отчистить чат")
)

func init() {
	StartKeyboard.Reply(
		StartKeyboard.Row(MissingBtn),
		StartKeyboard.Row(MyGroupsBtn, SettingsBtn),
		StartKeyboard.Row(AIBtn),
	)

	MyGroupsKeyboard.Inline(
		MyGroupsKeyboard.Row(refreshGroupsBtn),
	)

	SettingsKeyboard.Inline(
		SettingsKeyboard.Row(SetCookieBtn),
		SettingsKeyboard.Row(ChangeNotificationBtn),
	)

	RejectKeyboard.Reply(
		RejectKeyboard.Row(RejectActionBtn),
	)

	AIKeyboard.Reply(
		RejectKeyboard.Row(BackBtn, ClearHistoryBtn),
	)
}
