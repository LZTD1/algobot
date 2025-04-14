package keyboards

import tele "gopkg.in/telebot.v4"

func RefreshGroups() *tele.ReplyMarkup {
	refreshKb := &tele.ReplyMarkup{ResizeKeyboard: true}
	refresh := refreshKb.Data("Обновить группы", "refresh_groups")

	refreshKb.Inline(
		refreshKb.Row(refresh),
	)

	return refreshKb
}
