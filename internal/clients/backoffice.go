package clients

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBackofficeUrl = "https://backoffice.algoritmika.org"

type Backoffice struct {
	url      string
	settings BackofficeSetting
}
type BackofficeSetting struct {
	Retry        int
	Timeout      time.Duration
	RetryTimeout time.Duration
}

func NewBackoffice(url string, settings BackofficeSetting) *Backoffice {
	if strings.TrimSpace(url) == "" {
		url = defaultBackofficeUrl
	}
	return &Backoffice{url: url, settings: settings}
}

func (b Backoffice) GetKidsNamesByGroup(cookie, group string) (*GroupResponse, *ClientError) {
	req, err := b.createReq("GET", "/api/v2/group/student/index", cookie, map[string]string{
		"groupId": group,
		"expand":  "lastGroup",
	}, nil)
	if err != nil {
		return nil, GetError(500, err.Error())
	}

	data, reqErr := b.doReq(req)
	if reqErr != nil {
		return nil, reqErr
	}

	var response GroupResponse
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return nil, GetError(500, err.Error())
	}

	return &response, nil
}

func (b Backoffice) GetKidsStatsByGroup(cookie, group string) (*KidsStats, *ClientError) {
	req, err := b.createReq("GET", "/api/v1/stats/default/attendance", cookie, map[string]string{
		"group": group,
	}, nil)
	if err != nil {
		return nil, GetError(500, err.Error())
	}

	data, reqErr := b.doReq(req)
	if reqErr != nil {
		return nil, reqErr
	}

	var response KidsStats
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return nil, GetError(500, err.Error())
	}

	return &response, nil
}

func (b Backoffice) GetKidsMessages(cookie string) (*KidsMessages, *ClientError) {
	req, err := b.createReq("GET", "/api/v1/teacherComment/projects", cookie, map[string]string{
		"from":  "0",
		"limit": "30",
	}, nil)
	if err != nil {
		return nil, GetError(500, err.Error())
	}

	data, reqErr := b.doReq(req)
	if reqErr != nil {
		return nil, reqErr
	}

	var response KidsMessages
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return nil, GetError(500, err.Error())
	}

	return &response, nil
}

func (b Backoffice) GetAllGroupsByUser(cookie string) ([]AllGroupsUser, *ClientError) {
	req, err := b.createReq("GET", "/group", cookie, map[string]string{
		"GroupSearch[status][]": "active",
		"presetType":            "all",
		"_pjax":                 "#group-grid-pjax",
	}, nil)
	if err != nil {
		return nil, GetError(500, err.Error())
	}
	data, reqErr := b.doReq(req)
	if reqErr != nil {
		return nil, reqErr
	}

	res, err := parsedHtml(data.Body)
	if err != nil {
		return nil, GetError(500, err.Error())
	}
	return res, nil
}

func (b Backoffice) OpenLession(cookie, group, lession string) *ClientError {
	req, err := b.createReq("POST", "/api/v2/group/lesson/status", cookie, map[string]string{}, strings.NewReader("ajaxUrl=^%^2Fapi^%^2Fv2^%^2Fgroup^%^2Flesson^%^2Fstatus&btnClass=btn+btn-xs+btn-danger&status=10&lessonId="+lession+"&groupId="+group))
	if err != nil {
		return GetError(500, err.Error())
	}
	_, reqErr := b.doReq(req)
	if reqErr != nil {
		return reqErr
	}

	return nil
}

func (b Backoffice) CloseLession(cookie, group, lession string) *ClientError {
	req, err := b.createReq("POST", "/api/v2/group/lesson/status", cookie, map[string]string{}, strings.NewReader("ajaxUrl=^%^2Fapi^%^2Fv2^%^2Fgroup^%^2Flesson^%^2Fstatus&btnClass=btn+btn-xs+btn-danger&status=0&lessonId="+lession+"&groupId="+group))
	if err != nil {
		return GetError(500, err.Error())
	}
	_, reqErr := b.doReq(req)
	if reqErr != nil {
		return reqErr
	}

	return nil
}

func parsedHtml(body io.ReadCloser) ([]AllGroupsUser, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	var groups []AllGroupsUser

	doc.Find("tr.group-grid").Each(func(i int, row *goquery.Selection) {
		groupId := row.Find("td[data-col-seq='id']").First().Text()
		groupId = strings.TrimSpace(groupId)

		titleCell := row.Find("td[data-col-seq='title']").First()

		groupTitle := titleCell.Find("p").First().Text()
		groupTitle = strings.TrimSpace(groupTitle)

		groupTime := titleCell.Find("a").First().Text()
		groupTime = strings.TrimSpace(groupTime)

		nextLessonTime := row.Find("td[data-col-seq='nextLessonTime']").First().Text()
		nextLessonTime = strings.TrimSpace(nextLessonTime)

		groups = append(groups, AllGroupsUser{
			Title:       strings.ReplaceAll(groupTitle, "\u00A0", " "),
			GroupId:     strings.ReplaceAll(groupId, "\u00A0", " "),
			TimeLesson:  strings.ReplaceAll(nextLessonTime, "\u00A0", " "),
			RegularTime: strings.ReplaceAll(groupTime, "\u00A0", " "),
		})
	})

	return groups, nil
}
func (b Backoffice) doReq(req *http.Request) (*http.Response, *ClientError) {
	client := &http.Client{
		Timeout: b.settings.Timeout,
	}

	returnErr := &ClientError{}
	for i := 0; i < b.settings.Retry; i++ {
		resp, err := client.Do(req)
		if err != nil {
			returnErr = GetError(500, err.Error())
			time.Sleep(b.settings.RetryTimeout)
			continue
		}
		if resp.StatusCode >= 500 {
			returnErr = GetError(resp.StatusCode, getString(resp.Body))
			time.Sleep(b.settings.RetryTimeout)
			continue
		}
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, GetError(resp.StatusCode, getString(resp.Body))
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}
	}
	return nil, returnErr
}
func getString(body io.ReadCloser) string {
	all, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(all)
}
func (b Backoffice) createReq(method, uri, cookie string, params map[string]string, body io.Reader) (*http.Request, error) {
	reqUrl, _ := url.Parse(fmt.Sprintf("%s%s", b.url, uri))
	p := url.Values{}
	for key, val := range params {
		p.Add(key, val)
	}
	reqUrl.RawQuery = p.Encode()
	req, err := http.NewRequest(method, reqUrl.String(), body)
	if err != nil {
		return nil, GetError(500, err.Error())
	}
	req.Header.Add("Cookie", cookie)
	return req, nil
}
