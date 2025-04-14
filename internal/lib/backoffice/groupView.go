package backoffice

import (
	"algobot/internal/domain/backoffice"
	"encoding/json"
	"fmt"
)

func (bo *Backoffice) GroupView(groupID string, cookie string) (backoffice.GroupInfo, error) {
	const op = "backoffice.GetGroupView"

	req, err := bo.createReq("GET", "/api/v1/group/"+groupID, cookie, map[string]string{
		"expand": "venue,teacher,curator,branch",
	}, nil)
	if err != nil {
		return backoffice.GroupInfo{}, fmt.Errorf("%s err create req: %w", op, err)
	}

	data, err := bo.doReq(req)
	if err != nil {
		return backoffice.GroupInfo{}, fmt.Errorf("%s err doReq: %w", op, err)
	}

	var response backoffice.GroupInfo
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return backoffice.GroupInfo{}, fmt.Errorf("%s err decode json: %w", op, err)
	}

	return response, nil
}
