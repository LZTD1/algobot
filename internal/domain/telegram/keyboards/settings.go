package keyboards

import tele "gopkg.in/telebot.v4"

func Settings() *tele.ReplyMarkup {
	settingsKb := &tele.ReplyMarkup{ResizeKeyboard: true}
	setCookie := settingsKb.Data("Установить Cookie", "set_cookie")
	changeNotification := settingsKb.Data("Переключить уведомления", "change_notification")

	settingsKb.Inline(
		settingsKb.Row(setCookie),
		settingsKb.Row(changeNotification),
	)

	return settingsKb
}
