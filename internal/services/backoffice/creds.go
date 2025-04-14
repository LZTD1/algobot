package backoffice

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"log/slog"
)

func (bo *Backoffice) Creds(uid int64, groupID string, traceID interface{}) ([]models.Credential, error) {
	const op = "services.backoffice.Creds"
	log := bo.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	cookie, err := bo.cookieGetter.Cookies(uid)
	if err != nil {
		log.Warn("failed to get cookies", sl.Err(err))
		return nil, fmt.Errorf("%s failed to get cookies: %w", op, err)
	}
	group, err := bo.groupView.KidsNamesByGroup(groupID, cookie)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get KidsNamesByGroup: %w", op, err)
	}

	creds := make([]models.Credential, len(group.Data.Items))
	for i, item := range group.Data.Items {
		creds[i] = models.Credential{
			Fullname: item.FullName,
			Login:    item.Username,
			Password: item.Password,
		}
	}

	return creds, nil
}
