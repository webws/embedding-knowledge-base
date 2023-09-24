package util

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

func ProxyCheck(client *http.Client, proxy, urlTarget string, timeout int64) bool {
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cncl()
	if client == nil {
		client = &http.Client{}
	}
	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", urlTarget, nil)
	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
