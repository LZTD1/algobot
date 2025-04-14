package backoffice

import (
	"algobot/internal/domain/backoffice"
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/lib/mappers"
	"fmt"
	"log/slog"
	"strconv"
)

type GroupView interface {
	GroupView(groupID string, cookie string) (backoffice.GroupInfo, error)
	KidsNamesByGroup(groupId string, cookie string) (backoffice.NamesByGroup, error)
}

func (bo *Backoffice) GroupView(uid int64, groupID string, traceID interface{}) (models.GroupView, error) {
	const op = "services.backoffice.GetGroupView"
	log := bo.log.With(
		slog.String("op", op),
		slog.Any("traceID", traceID),
	)

	cookie, err := bo.cookieGetter.Cookies(uid)
	if err != nil {
		log.Warn("failed to get cookies", sl.Err(err))
		return models.GroupView{}, fmt.Errorf("%s failed to get cookies: %w", op, err)
	}

	grView, err := bo.groupView.GroupView(groupID, cookie)
	if err != nil {
		log.Warn("failed to group view", sl.Err(err))
		return models.GroupView{}, fmt.Errorf("%s failed to group view: %w", op, err)
	}
	kidsNames, err := bo.groupView.KidsNamesByGroup(groupID, cookie)
	if err != nil {
		log.Warn("failed to get kids names", sl.Err(err))
		return models.GroupView{}, fmt.Errorf("%s failed to get kids names: %w", op, err)
	}

	return mapResp(grView, kidsNames, groupID), nil
}

func mapResp(info backoffice.GroupInfo, names backoffice.NamesByGroup, groupID string) models.GroupView {
	m := models.GroupView{}
	m.GroupID = info.Data.ID
	m.GroupTitle = info.Data.Title
	m.GroupContent = info.Data.Content
	m.NextLessonTime = info.Data.NextLessonTime
	m.LessonsPassed = info.Data.LessonsPassed
	m.LessonsTotal = info.Data.LessonsTotal

	for _, item := range names.Data.Items {
		if strconv.Itoa(item.LastGroup.ID) == groupID && item.LastGroup.Status == 0 {
			m.ActiveKids = append(m.ActiveKids, mappers.MapKid(item))
		} else {
			m.NotActiveKids = append(m.NotActiveKids, mappers.MapKid(item))
		}
	}

	return m
}
