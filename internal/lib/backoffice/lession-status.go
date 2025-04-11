package backoffice

import (
	"fmt"
	"net/url"
	"strings"
)

func (bo *Backoffice) OpenLesson(cookie, group, lession string) error {
	const op = "backoffice.OpenLesson"

	params := url.Values{}
	params.Add("ajaxUrl", "/api/v2/group/lesson/status")
	params.Add("btnClass", "btn btn-xs btn-danger")
	params.Add("status", "10")
	params.Add("lessonId", lession)
	params.Add("groupId", group)
	query := params.Encode()

	req, err := bo.createReq("POST", "/api/v2/group/lesson/status", cookie, map[string]string{}, strings.NewReader(query))
	if err != nil {
		return fmt.Errorf("%s err createReq: %w", op, err)
	}
	_, reqErr := bo.doReq(req)
	if reqErr != nil {
		return fmt.Errorf("%s err doReq: %w", op, err)
	}

	return nil
}

func (bo *Backoffice) CloseLesson(cookie, group, lession string) error {
	const op = "backoffice.CloseLesson"

	params := url.Values{}
	params.Add("ajaxUrl", "/api/v2/group/lesson/status")
	params.Add("btnClass", "btn btn-xs btn-danger")
	params.Add("status", "0")
	params.Add("lessonId", lession)
	params.Add("groupId", group)
	query := params.Encode()

	req, err := bo.createReq("POST", "/api/v2/group/lesson/status", cookie, map[string]string{}, strings.NewReader(query))
	if err != nil {
		return fmt.Errorf("%s err createReq: %w", op, err)
	}

	_, reqErr := bo.doReq(req)
	if reqErr != nil {
		return fmt.Errorf("%s err doReq: %w", op, err)
	}

	return nil
}
