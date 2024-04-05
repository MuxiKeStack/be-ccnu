package service

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func (c *ccnuService) Login(ctx context.Context, studentId string, password string) (bool, error) {
	client, err := c.loginClient(ctx, studentId, password)
	return client != nil, err
}

func (c *ccnuService) client() *http.Client {
	j, _ := cookiejar.New(&cookiejar.Options{})
	return &http.Client{
		Timeout: c.timeout,
		Jar:     j,
	}
}

func (c *ccnuService) loginClient(ctx context.Context, studentId string, password string) (*http.Client, error) {
	params, err := c.makeAccountPreflightRequest()
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("username", studentId)
	v.Set("password", password)
	v.Set("lt", params.lt)
	v.Set("execution", params.execution)
	v.Set("_eventId", params._eventId)
	v.Set("submit", params.submit)

	request, err := http.NewRequest("POST", "https://account.ccnu.edu.cn/cas/login;jsessionid="+params.JSESSIONID, strings.NewReader(v.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")

	client := c.client()
	resp, err := client.Do(request)
	if err != nil || len(resp.Header.Get("Set-Cookie")) == 0 {
		return nil, err
	}
	return client, err
}
