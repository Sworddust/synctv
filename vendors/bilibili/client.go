package bilibili

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/synctv-org/synctv/utils"
)

type Client struct {
	httpClient *http.Client
	cookies    []*http.Cookie
	buvid3     *http.Cookie
	once       sync.Once
	ctx        context.Context
}

type ClientConfig func(*Client)

func WithHttpClient(httpClient *http.Client) ClientConfig {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithContext(ctx context.Context) ClientConfig {
	return func(c *Client) {
		c.ctx = ctx
	}
}

func NewClient(cookies []*http.Cookie, conf ...ClientConfig) (*Client, error) {
	cli := &Client{
		httpClient: http.DefaultClient,
		cookies:    cookies,
		ctx:        context.Background(),
	}
	for _, v := range conf {
		v(cli)
	}
	return cli, nil
}

func (c *Client) InitBuvid3() (err error) {
	c.once.Do(func() {
		c.buvid3, err = newBuvid3(c.ctx)
	})
	return
}

func newBuvid3(ctx context.Context) (*http.Cookie, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.bilibili.com/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", utils.UA)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	for _, c := range resp.Cookies() {
		if c.Name == "buvid3" {
			return c, nil
		}
	}
	return nil, errors.New("no buvid3 cookie")
}

func (c *Client) SetCookies(cookies []*http.Cookie) {
	c.cookies = cookies
}

type RequestConfig struct {
	wbi bool
}

func defaultRequestConfig() *RequestConfig {
	return &RequestConfig{
		wbi: true,
	}
}

type RequestOption func(*RequestConfig)

func WithoutWbi() RequestOption {
	return func(c *RequestConfig) {
		c.wbi = false
	}
}

func (c *Client) NewRequest(method, url string, body io.Reader, conf ...RequestOption) (req *http.Request, err error) {
	config := defaultRequestConfig()
	for _, v := range conf {
		v(config)
	}
	if config.wbi {
		url, err = signAndGenerateURL(url)
		if err != nil {
			return nil, err
		}
	}
	req, err = http.NewRequestWithContext(c.ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if c.buvid3 != nil {
		req.AddCookie(c.buvid3)
	}
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}
	req.Header.Set("User-Agent", utils.UA)
	req.Header.Set("Referer", "https://www.bilibili.com")
	return req, nil
}
