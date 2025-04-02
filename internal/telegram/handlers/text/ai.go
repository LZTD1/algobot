package text

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/fsm"
	"algobot/internal/lib/logger/sl"
	"gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

type AIInfo struct {
	TextModel  string
	ImageModel string
}

type AIInformer interface {
	GetAIInfo() (AIInfo, error)
}

type AIStater interface {
	SetState(uid int64, state fsm.State)
}

func NewAI(ai AIInformer, log *slog.Logger, stater AIStater) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "text.NewAI"

		log = log.With(
			slog.String("op", op),
			slog.Any("trace_id", ctx.Get("trace_id")),
		)
		uid := ctx.Sender().ID

		info, err := ai.GetAIInfo()
		if err != nil {
			log.Warn("error while GetAIInfo", sl.Err(err))
			return ctx.Send("Упс, AI сейчас не работает!")
		}

		stater.SetState(uid, fsm.ChattingAI)
		return ctx.Send(GetAIMessage(info), keyboards.RejectKeyboard())
	}
}

func GetAIMessage(info AIInfo) string {
	sb := strings.Builder{}
	sb.WriteString("Информация о включенных моделях:\n\n")
	sb.WriteString("Текст: ")
	sb.WriteString(info.TextModel)
	sb.WriteString("\nИзображение: ")
	sb.WriteString(info.ImageModel)
	sb.WriteString("\n\n")
	sb.WriteString("/reset - отчистить память модели")
	sb.WriteString("\n/image promt - сгенерировать изображение")
	sb.WriteString("\nДля текстового запроса - просто напиши в чат")
	return sb.String()
}
