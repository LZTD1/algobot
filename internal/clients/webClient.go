package clients

type WebClient interface {
	// GetKidsNamesByGroup получить всех детей в группе
	GetKidsNamesByGroup(cookie, group string)
	// GetKidsStatsByGroup получить статистику посещения детей в группе
	GetKidsStatsByGroup(cookie, group string)
	// OpenLession открыть лекцию с идентификатором {lession}
	OpenLession(cookie, group, lession string)
	// CloseLession закрыть лекцию с идентификатором {lession}
	CloseLession(cookie, group, lession string)
	// GetKidsMessages получить новые сообщения детей на платформе
	GetKidsMessages(cookie string)
	// GetAllGroupsByUser получить все группы
	GetAllGroupsByUser(cookie string)
}
