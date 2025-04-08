package backoffice

import (
	"algobot/internal/config"
	"algobot/internal/lib/logger/handlers/slogpretty"
	"algobot/internal/lib/logger/sl"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrBadCode = errors.New("bad code")
)

type Option func(*Backoffice)

type Backoffice struct {
	url    string
	client *http.Client
	log    *slog.Logger
	cfg    *config.Backoffice
}

func NewBackoffice(cfg *config.Backoffice, fn ...Option) *Backoffice {
	bo := &Backoffice{
		cfg: cfg,
		url: "https://backoffice.algoritmika.org",
		client: &http.Client{
			Timeout: cfg.ResponseTimeout,
		},
		log: slog.New(slogpretty.NewHandler(&slog.HandlerOptions{Level: slog.LevelDebug})),
	}
	for _, o := range fn {
		o(bo)
	}

	return bo
}

func WithURL(url string) func(*Backoffice) {
	return func(bo *Backoffice) {
		bo.url = url
	}
}
func WithLogger(log *slog.Logger) func(*Backoffice) {
	return func(bo *Backoffice) {
		bo.log = log
	}
}

func (bo *Backoffice) doReq(req *http.Request) (*http.Response, error) {
	const op = "backoffice.doReq"
	log := bo.log.With(
		slog.String("op", op),
	)

	var err error
	var resp *http.Response

	for i := 0; i < bo.cfg.Retries; i++ {
		resp, err = bo.client.Do(req)
		if err != nil {
			log.Debug("error while req bo", sl.Err(err))
			time.Sleep(bo.cfg.RetriesTimeout)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Debug("received not 200 OK", slog.Int("status", resp.StatusCode))
			err = ErrBadCode
			time.Sleep(bo.cfg.RetriesTimeout)
			continue
		}
		return resp, nil
	}

	if err != nil {
		log.Warn("error while req bo", sl.Err(err))
		return nil, fmt.Errorf("%s error while trying send request: %w", op, err)
	}

	return resp, nil
}
func (bo *Backoffice) createReq(method, uri, cookie string, params map[string]string, body io.Reader) (*http.Request, error) {
	const op = "backoffice.createReq"

	reqUrl, _ := url.Parse(fmt.Sprintf("%s%s", bo.url, uri))

	p := url.Values{}
	for key, val := range params {
		p.Add(key, val)
	}
	reqUrl.RawQuery = p.Encode()

	req, err := http.NewRequest(method, reqUrl.String(), body)
	if err != nil {
		return nil, fmt.Errorf("%s error while creating req: %w", op, err)
	}

	req.Header.Add("Cookie", cookie)
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	return req, nil
}
