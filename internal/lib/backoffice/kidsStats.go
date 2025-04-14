package backoffice

import (
	"algobot/internal/domain/backoffice"
	"encoding/json"
	"fmt"
	"strconv"
)

func (bo *Backoffice) KidsStats(cookie string, groupID int) (backoffice.KidsStats, error) {
	const op = "backoffice.KidsStats"

	req, err := bo.createReq("GET", "/api/v1/stats/default/attendance", cookie, map[string]string{
		"group": strconv.Itoa(groupID),
	}, nil)
	if err != nil {
		return backoffice.KidsStats{}, fmt.Errorf("%s err createReq: %w", op, err)
	}

	data, reqErr := bo.doReq(req)
	if reqErr != nil {
		return backoffice.KidsStats{}, fmt.Errorf("%s err doReq: %w", op, err)
	}

	var response backoffice.KidsStats
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return backoffice.KidsStats{}, fmt.Errorf("%s err while docding: %w", op, err)
	}

	return response, nil
}
