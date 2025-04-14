package backoffice

import (
	"algobot/internal/domain/backoffice"
	"encoding/json"
	"fmt"
)

func (bo *Backoffice) KidsNamesByGroup(groupId string, cookie string) (backoffice.NamesByGroup, error) {
	const op = "backoffice.KidsNamesByGroup"

	req, err := bo.createReq("GET", "/api/v2/group/student/index", cookie, map[string]string{
		"groupId": groupId,
		"expand":  "lastGroup, groups",
	}, nil)
	if err != nil {
		return backoffice.NamesByGroup{}, fmt.Errorf("%s cant create req: %w", op, err)
	}

	data, reqErr := bo.doReq(req)
	if reqErr != nil {
		return backoffice.NamesByGroup{}, fmt.Errorf("%s cant do req: %w", op, err)
	}

	var response backoffice.NamesByGroup
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return backoffice.NamesByGroup{}, fmt.Errorf("%s cant decode res: %w", op, err)
	}

	return response, nil
}
