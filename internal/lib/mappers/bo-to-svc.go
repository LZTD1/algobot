package mappers

import (
	"algobot/internal/domain/backoffice"
	"algobot/internal/domain/models"
)

func MapKid(item backoffice.Student) models.GroupKid {
	m := models.GroupKid{}
	m.ID = item.ID
	m.FullName = item.FullName
	m.LastGroup = MapGroup(item.LastGroup)

	return m
}

func MapGroup(group backoffice.Group) models.KidGroup {
	m := models.KidGroup{}
	m.ID = group.ID
	m.StartTime = group.StartTime
	m.EndTime = group.EndTime

	return m
}

func MapKidView(view backoffice.KidView) models.KidView {
	m := models.KidView{}
	m.Kid.FullName = view.Data.FullName
	m.Kid.ParentName = view.Data.ParentName
	m.Kid.Email = view.Data.Email
	m.Kid.Phone = view.Data.Phone
	m.Kid.Age = view.Data.Age
	m.Kid.BirthDate = view.Data.BirthDate
	m.Kid.Username = view.Data.Username
	m.Kid.Password = view.Data.Password

	m.Kid.Groups = make([]models.KidViewGroup, len(view.Data.Groups))
	for i, group := range view.Data.Groups {
		m.Kid.Groups[i] = models.KidViewGroup{
			ID:        group.ID,
			Title:     group.Title,
			Content:   group.Content,
			Status:    group.Status,
			StartTime: group.StartTime,
			EndTime:   group.EndTime,
		}
	}

	return m
}

func MapGroups(groups []backoffice.Group) []models.KidViewGroup {
	mappedGroups := make([]models.KidViewGroup, len(groups))
	for i, group := range groups {
		mappedGroups[i] = models.KidViewGroup{
			ID:        group.ID,
			Title:     group.Title,
			Content:   group.Content,
			Status:    group.Status,
			StartTime: group.StartTime,
			EndTime:   group.EndTime,
		}
	}
	return mappedGroups
}
