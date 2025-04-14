package backoffice

import (
	backoffice2 "algobot/internal/domain/backoffice"
	"encoding/json"
	"fmt"
)

func (bo *Backoffice) KidView(kidID string, cookie string) (backoffice2.KidView, error) {
	const op = "backoffice.KidView"

	req, err := bo.createReq("GET", "/api/v2/student/default/view/"+kidID, cookie, map[string]string{
		"expand": "groups",
	}, nil)

	if err != nil {
		return backoffice2.KidView{}, fmt.Errorf("%s createReq err: %w", op, err)
	}
	data, err := bo.doReq(req)
	if err != nil {
		return backoffice2.KidView{}, fmt.Errorf("%s do req err: %w", op, err)
	}

	var response backoffice2.KidView
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return backoffice2.KidView{}, fmt.Errorf("%s decode res err: %w", op, err)
	}

	return response, nil
}
