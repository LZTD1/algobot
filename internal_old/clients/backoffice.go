package clients

import (
	appError "algobot/internal_old/error"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

func (b Backoffice) GetKidInfo(cookie string, kidID string) (*FullKidInfo, error) {
	req, err := b.createReq("GET", "/api/v2/student/default/view/"+kidID, cookie, map[string]string{
		"expand": "groups",
	}, nil)

	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetKidInfo(%s, %s) : %w", cookie, kidID, err)
	}
	data, err := b.doReq(req)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetKidInfo(%s, %s) : %w", cookie, kidID, err)
	}

	var response FullKidInfo
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetKidInfo(%s, %s) : %w", cookie, kidID, err)
	}

	return &response, nil
}

func (b Backoffice) GetGroupInfo(cookie string, group string) (*FullGroupInfo, error) {
	req, err := b.createReq("GET", "/api/v1/group/"+group, cookie, map[string]string{
		"expand": "venue,teacher,curator,branch",
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetGroupInfo(%s, %s) : %w", cookie, group, err)
	}
	data, err := b.doReq(req)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetGroupInfo(%s, %s) : %w", cookie, group, err)
	}

	var response FullGroupInfo
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetGroupInfo(%s, %s) : %w", cookie, group, err)
	}

	return &response, nil
}

func (b Backoffice) GetKidsNamesByGroup(cookie string, group int) (*GroupResponse, error) {

	req, err := b.createReq("GET", "/api/v2/group/student/index", cookie, map[string]string{
		"groupId": strconv.Itoa(group),
		"expand":  "lastGroup, groups",
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetKidsNamesByGroup(%s, %s) : %w", cookie, group, err)
	}

	data, reqErr := b.doReq(req)
	if reqErr != nil {
		return nil, fmt.Errorf("Backoffice.GetKidsNamesByGroup(%s, %s) : %w", cookie, group, reqErr)
	}

	var response GroupResponse
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetKidsNamesByGroup(%s, %s) : %w", cookie, group, err)
	}

	return &response, nil
}

func (b Backoffice) GetKidsStatsByGroup(cookie, group string) (*KidsStats, error) {
	req, err := b.createReq("GET", "/api/v1/stats/default/attendance", cookie, map[string]string{
		"group": group,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetKidsStatsByGroup(%s, %s) : %w", cookie, group, err)
	}

	data, reqErr := b.doReq(req)
	if reqErr != nil {
		return nil, fmt.Errorf("Backoffice.GetKidsStatsByGroup(%s, %s) : %w", cookie, group, reqErr)
	}

	var response KidsStats
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetKidsStatsByGroup(%s, %s) : %w", cookie, group, err)
	}

	return &response, nil
}

func (b Backoffice) GetKidsMessages(cookie string) (*KidsMessages, error) {
	req, err := b.createReq("GET", "/api/v1/teacherComment/projects", cookie, map[string]string{
		"from":  "0",
		"limit": "30",
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetKidsMessages(%s) : %w", cookie, err)
	}

	data, reqErr := b.doReq(req)
	if reqErr != nil {
		return nil, fmt.Errorf("Backoffice.GetKidsMessages(%s) : %w", cookie, reqErr)
	}

	var response KidsMessages
	err = json.NewDecoder(data.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetKidsMessages(%s) : %w", cookie, err)
	}

	return &response, nil
}

func (b Backoffice) GetAllGroupsByUser(cookie string) ([]AllGroupsUser, error) {
	req, err := b.createReq("GET", "/group", cookie, map[string]string{
		"GroupSearch[status][]": "active",
		"presetType":            "all",
		"_pjax":                 "#group-grid-pjax",
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetAllGroupsByUser(%s) : %w", cookie, err)
	}
	data, reqErr := b.doReq(req)
	if reqErr != nil {
		return nil, fmt.Errorf("Backoffice.GetAllGroupsByUser(%s) : %w", cookie, reqErr)
	}

	res, err := parsedHtml(data.Body)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.GetAllGroupsByUser(%s) : %w", cookie, err)
	}
	return res, nil
}

func (b Backoffice) OpenLession(cookie, group, lession string) error {
	params := url.Values{}
	params.Add("ajaxUrl", "/api/v2/group/lesson/status")
	params.Add("btnClass", "btn btn-xs btn-danger")
	params.Add("status", "10")
	params.Add("lessonId", lession)
	params.Add("groupId", group)
	query := params.Encode()

	req, err := b.createReq("POST", "/api/v2/group/lesson/status", cookie, map[string]string{}, strings.NewReader(query))
	if err != nil {
		return fmt.Errorf("Backoffice.OpenLession(%s, %s, %s) : %w", cookie, group, lession, err)
	}
	_, reqErr := b.doReq(req)
	if reqErr != nil {
		return fmt.Errorf("Backoffice.OpenLession(%s, %s, %s) : %w", cookie, group, lession, reqErr)
	}

	return nil
}

func (b Backoffice) CloseLession(cookie, group, lession string) error {
	params := url.Values{}
	params.Add("ajaxUrl", "/api/v2/group/lesson/status")
	params.Add("btnClass", "btn btn-xs btn-danger")
	params.Add("status", "0")
	params.Add("lessonId", lession)
	params.Add("groupId", group)
	query := params.Encode()

	req, err := b.createReq("POST", "/api/v2/group/lesson/status", cookie, map[string]string{}, strings.NewReader(query))
	if err != nil {
		return fmt.Errorf("Backoffice.CloseLession(%s, %s, %s) : %w", cookie, group, lession, err)
	}

	_, reqErr := b.doReq(req)
	if reqErr != nil {
		return fmt.Errorf("Backoffice.CloseLession(%s, %s, %s) : %w", cookie, group, lession, reqErr)
	}

	return nil
}

func parsedHtml(body io.ReadCloser) ([]AllGroupsUser, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("Backoffice.parsedHtml() : %w", err)
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

		if nextLessonTime != "" {
			groups = append(groups, AllGroupsUser{
				Title:       strings.ReplaceAll(groupTitle, "\u00A0", " "),
				GroupId:     strings.ReplaceAll(groupId, "\u00A0", " "),
				TimeLesson:  strings.ReplaceAll(nextLessonTime, "\u00A0", " "),
				RegularTime: strings.ReplaceAll(groupTime, "\u00A0", " "),
			})
		}
	})

	return groups, nil
}
func (b Backoffice) doReq(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: b.settings.Timeout,
	}

	var returnErr error

	for i := 0; i < b.settings.Retry; i++ {
		resp, err := client.Do(req)
		if err != nil {
			returnErr = fmt.Errorf("Backoffice.doReq() : %w", err)
			time.Sleep(b.settings.RetryTimeout)
			continue
		}
		if resp.StatusCode >= 500 {
			returnErr = fmt.Errorf("Backoffice.doReq() : %w", errors.New(resp.Status+" "+getString(resp.Body)))
			time.Sleep(b.settings.RetryTimeout)
			continue
		}
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, fmt.Errorf("Backoffice.doReq() : %w : %w", appError.ErrNotFound, getErrorByCode(resp.Status, resp.Body))
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}
	}
	return nil, returnErr
}

func getErrorByCode(status string, body io.Reader) error {
	return errors.New(status + " " + getString(body))
}

func getString(body io.Reader) string {
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
		return nil, fmt.Errorf("Backoffice.createReq() : %w", err)
	}
	req.Header.Add("Cookie", cookie)

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	return req, nil
}
