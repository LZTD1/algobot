package backoffice

import (
	"algobot/internal/domain/models"
	"encoding/json"
	"fmt"
)

func (bo *Backoffice) GroupView(uid int64, groupID string, cookie string) (models.GroupView, error) {
	const op = "backoffice.GetGroupView"

	req, err := bo.createReq("GET", "/api/v1/group/"+groupID, cookie, map[string]string{
		"expand": "venue,teacher,curator,branch",
	}, nil)
	if err != nil {
		return models.GroupView{}, fmt.Errorf("%s err create req: %w", op, err)
	}

	data, err := bo.doReq(req)
	if err != nil {
		return models.GroupView{}, fmt.Errorf("%s err doReq: %w", op, err)
	}

	var response models.GroupView
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return models.GroupView{}, fmt.Errorf("%s err decode json: %w", op, err)
	}

	return models.GroupView{}, nil
}
