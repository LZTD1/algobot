package backoffice

import "log/slog"

type CookieGetter interface {
	Cookies(uid int64) (string, error)
}

type Backoffice struct {
	log            *slog.Logger
	cookieGetter   CookieGetter
	groupView      GroupView
	kidViewer      KidViewer
	lessonStatuser LessonStatuser
}

func NewBackoffice(log *slog.Logger, cookieGetter CookieGetter, groupView GroupView, kidViewer KidViewer, lessonStatus LessonStatuser) *Backoffice {
	return &Backoffice{log: log, cookieGetter: cookieGetter, groupView: groupView, kidViewer: kidViewer, lessonStatuser: lessonStatus}
}
