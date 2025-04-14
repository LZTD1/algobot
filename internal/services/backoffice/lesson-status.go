package backoffice

import (
	"algobot/internal/lib/logger/sl"
	"fmt"
	"log/slog"
)

type LessonStatus int

const (
	CloseLesson LessonStatus = iota
	OpenLesson
)

type LessonStatuser interface {
	OpenLesson(cookie, group, lession string) error
	CloseLesson(cookie, group, lession string) error
}

func (bo *Backoffice) SetLessonStatus(uid int64, groupID string, lessonID string, status LessonStatus, traceID interface{}) error {
	const op = "services.backoffice.GetGroupView"
	log := bo.log.With(
		slog.String("op", op),
		slog.Any("traceID", traceID),
	)

	cookie, err := bo.cookieGetter.Cookies(uid)
	if err != nil {
		log.Warn("failed to get cookies", sl.Err(err))
		return fmt.Errorf("%s failed to get cookies: %w", op, err)
	}

	switch status {
	case CloseLesson:
		if err := bo.lessonStatuser.CloseLesson(cookie, groupID, lessonID); err != nil {
			return fmt.Errorf("%s error while CloseLesson : %w", op, err)
		}
	case OpenLesson:
		if err := bo.lessonStatuser.OpenLesson(cookie, groupID, lessonID); err != nil {
			return fmt.Errorf("%s error while OpenLesson : %w", op, err)
		}
	default:
		return fmt.Errorf("%s invalid lesson status: %d", op, status)
	}

	return nil
}
