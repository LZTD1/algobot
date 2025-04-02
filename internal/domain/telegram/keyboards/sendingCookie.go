package keyboards

import tele "gopkg.in/telebot.v4"

func RejectKeyboard() *tele.ReplyMarkup {
	rejectKb := &tele.ReplyMarkup{ResizeKeyboard: true}
	reject := rejectKb.Text("⬅️ Назад")

	rejectKb.Reply(
		rejectKb.Row(reject),
	)

	return rejectKb
}
