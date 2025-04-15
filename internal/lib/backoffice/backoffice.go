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
	ErrBadCode   = errors.New("bad code")
	Err4xxStatus = errors.New("not found")
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
		slog.String("req", req.URL.String()),
	)

	var errChain error
	var resp *http.Response

	for i := 0; i < bo.cfg.Retries; i++ {
		resp, err := bo.client.Do(req)
		if err != nil {
			log.Debug("error while req bo", sl.Err(err))
			time.Sleep(bo.cfg.RetriesTimeout)
			errChain = fmt.Errorf("%w err do req: %w", errChain, err)
			continue
		}
		if resp.StatusCode >= 400 && resp.StatusCode <= 499 {
			log.Debug("received 4xx", slog.Int("status", resp.StatusCode))
			return nil, fmt.Errorf("%s : %w", op, Err4xxStatus)
		}
		if resp.StatusCode != http.StatusOK {
			log.Debug("received not 200 OK", slog.Int("status", resp.StatusCode))
			errChain = fmt.Errorf("%w bad code %d : %w", errChain, resp.StatusCode, ErrBadCode)
			continue
		}
		return resp, nil
	}

	if errChain != nil {
		log.Warn("error while req bo", sl.Err(errChain))
		return nil, fmt.Errorf("%s error while trying send request: %w", op, errChain)
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
