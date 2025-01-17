package clients

import (
	"fmt"
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
	uri, _ := url.Parse(fmt.Sprintf("%s%s", b.url, "/api/v2/group/student/index"))
	p := url.Values{}
	p.Add("groupId", group)
	p.Add("expand", "lastGroup")
	uri.RawQuery = p.Encode()

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, GetError(500, err.Error())
	}
	req.Header.Add("Cookie", cookie)
	for i := 0; i < b.settings.Retry; i++ {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, GetError(500, err.Error())
		}
		if resp.StatusCode != 200 {
			return nil, GetError(resp.StatusCode, getString(resp.Body))
		}
		if resp.StatusCode >= 500 {
			time.Sleep(b.settings.RetryTimeout)
		}
	}
	вуа
	return nil, nil
}

func (b Backoffice) GetKidsStatsByGroup(cookie, group string) {
	//TODO implement me
	panic("implement me")
}

func (b Backoffice) OpenLession(cookie, group, lession string) {
	//TODO implement me
	panic("implement me")
}

func (b Backoffice) CloseLession(cookie, group, lession string) {
	//TODO implement me
	panic("implement me")
}

func (b Backoffice) GetKidsMessages(cookie string) {
	//TODO implement me
	panic("implement me")
}

func (b Backoffice) GetAllGroupsByUser(cookie string) {
	//TODO implement me
	panic("implement me")
}

func getString(body io.ReadCloser) string {
	all, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(all)
}
