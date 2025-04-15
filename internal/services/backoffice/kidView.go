package backoffice

import (
	backoffice2 "algobot/internal/domain/backoffice"
	"algobot/internal/domain/models"
	"algobot/internal/lib/backoffice"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/lib/mappers"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
)

type KidViewer interface {
	KidView(kidID string, cookie string) (backoffice2.KidView, error)
	KidsNamesByGroup(groupId string, cookie string) (backoffice2.NamesByGroup, error)
}

func (bo *Backoffice) KidView(uid int64, kidID string, groupId string, traceID interface{}) (models.KidView, error) {
	const op = "services.backoffice.KidView"
	log := bo.log.With(
		slog.String("op", op),
		slog.Any("traceID", traceID),
	)

	cookie, err := bo.cookieGetter.Cookies(uid)
	if err != nil {
		log.Warn("failed to get cookies", sl.Err(err))
		return models.KidView{}, fmt.Errorf("%s failed to get cookies: %w", op, err)
	}

	view, err := bo.kidViewer.KidView(kidID, cookie)
	if err != nil {
		if errors.Is(err, backoffice.Err4xxStatus) { // TODO : maybe refactor into single request
			info, err := bo.kidViewer.KidsNamesByGroup(groupId, cookie)
			if err != nil {
				log.Warn("failed to get kids names by group", sl.Err(err))
				return models.KidView{}, fmt.Errorf("%s failed to get kids names by group: %w", op, err)
			}
			for _, item := range info.Data.Items {
				if strconv.Itoa(item.ID) == kidID {
					return models.KidView{
						Extra: models.NotAccessible,
						Kid: models.Kid{
							FullName:   item.FullName,
							ParentName: item.ParentName,
							Email:      item.Email,
							Phone:      item.Phone,
							Age:        item.Age,
							BirthDate:  item.BirthDate,
							Username:   item.Username,
							Password:   item.Password,
							Groups:     mappers.MapGroups(item.Groups),
						},
					}, nil
				}
			}
		}
		log.Warn("failed to kid view", sl.Err(err))
		return models.KidView{}, fmt.Errorf("%s failed to kid view: %w", op, err)
	}

	return mappers.MapKidView(view), nil
}
