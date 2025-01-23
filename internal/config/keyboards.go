package config

import tele "gopkg.in/telebot.v4"

var (
	StartKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}
	missing       = StartKeyboard.Text("Получить отсутсвующих")
	myGroups      = StartKeyboard.Text("Мои группы")
	settings      = StartKeyboard.Text("Настройки")

	MyGroupsKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}
	refreshGroups    = MyGroupsKeyboard.Data("Обновить группы", "refresh_groups")

	SettingsKeyboard   = &tele.ReplyMarkup{ResizeKeyboard: true}
	setCookie          = SettingsKeyboard.Data("Установить Cookie", "set_cookie")
	changeNotification = SettingsKeyboard.Data("Переключить уведомления", "change_notification")

	RejectKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}
	rejectAction   = RejectKeyboard.Text("Отменить действие")
)

func init() {
	StartKeyboard.Reply(
		StartKeyboard.Row(missing),
		StartKeyboard.Row(myGroups, settings),
	)

	MyGroupsKeyboard.Inline(
		MyGroupsKeyboard.Row(refreshGroups),
	)

	SettingsKeyboard.Inline(
		SettingsKeyboard.Row(setCookie),
		SettingsKeyboard.Row(changeNotification),
	)

	RejectKeyboard.Reply(
		RejectKeyboard.Row(rejectAction),
	)
}
