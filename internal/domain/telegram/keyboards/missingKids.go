package keyboards

import (
	"fmt"
	"gopkg.in/telebot.v4"
)

func MissingKids(groupID, lessonID int) *telebot.ReplyMarkup {
	markup := telebot.ReplyMarkup{ResizeKeyboard: true}
	markup.Inline(
		markup.Row(
			markup.Data("Закрыть лекцию", fmt.Sprintf("close_lesson_%d_%d", groupID, lessonID)),
			markup.Data("Открыть лекцию", fmt.Sprintf("open_lesson_%d_%d", groupID, lessonID)),
		),
		markup.Row(markup.Data("Получить аккаунты", fmt.Sprintf("get_creds_%d", groupID))),
	)

	return &markup
}
