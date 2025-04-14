package backoffice

import (
	"algobot/internal/domain/backoffice"
	"encoding/json"
	"fmt"
)

func (bo *Backoffice) KidsMessages(cookie string) (backoffice.KidsMessages, error) {
	const op = "backoffice.KidsMessages"

	req, err := bo.createReq("GET", "/api/v1/teacherComment/projects", cookie, map[string]string{
		"from":  "0",
		"limit": "30",
	}, nil)
	if err != nil {
		return backoffice.KidsMessages{}, fmt.Errorf("%s err create req: %w", op, err)
	}

	data, reqErr := bo.doReq(req)
	if reqErr != nil {
		return backoffice.KidsMessages{}, fmt.Errorf("%s err doReq: %w", op, err)
	}

	var response backoffice.KidsMessages
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return backoffice.KidsMessages{}, fmt.Errorf("%s err decode json: %w", op, err)
	}

	return response, nil
}
