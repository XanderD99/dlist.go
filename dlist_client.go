package dlist

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Config struct {
	Version  string `required:"true"`
	BotID    string
	APIToken string `required:"true"`

	RequestTimeout time.Duration `default:"5s"`
}

type DlistClient struct {
	BaseUrl  string
	BotID    string
	APIToken string

	http *http.Client
}

func New(cfg *Config) *DlistClient {
	return &DlistClient{
		BaseUrl:  fmt.Sprintf("https://api.discordlist.gg/%s", cfg.Version),
		BotID:    cfg.BotID,
		APIToken: cfg.APIToken,
		http: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
	}
}

func (d *DlistClient) request(ctx context.Context, method, url string, body []byte, headers http.Header) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", d.BaseUrl, url), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", d.APIToken))

	for h, val := range headers {
		for _, item := range val {
			req.Header.Add(h, item)
		}
	}

	res, err := d.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("received unexpected status code: %v", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading body: %v", err)
	}

	return resBody, nil
}

func (d *DlistClient) get(ctx context.Context, url string, headers http.Header) ([]byte, error) {
	return d.request(ctx, http.MethodGet, url, nil, http.Header{})
}

func (d *DlistClient) post(ctx context.Context, url string, body []byte, headers http.Header) ([]byte, error) {
	return d.request(ctx, http.MethodPost, url, body, headers)
}

func (d *DlistClient) put(ctx context.Context, url string, body []byte, headers http.Header) ([]byte, error) {
	return d.request(ctx, http.MethodPut, url, body, headers)
}

func (d *DlistClient) PostGuilds(count int) ([]byte, error) {
	return d.put(context.Background(), fmt.Sprintf("/bots/%s/guilds?count=%d", d.BotID, count), nil, http.Header{})
}

func (d *DlistClient) GetBot(ID string) ([]byte, error) {
	return d.get(context.Background(), fmt.Sprintf("/bots/%s", ID), http.Header{})
}
