package proxy

import (
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// CreateProxy return an http.Client with a proxy
func CreateProxy(proxylist []string, user, pass string) *http.Client {
	client := &http.Client{Timeout: 30 * time.Second}
	proxyUrl := &url.URL{Scheme: "http", Host: proxylist[rand.Intn(len(proxylist))]}
	if user != "" && pass != "" {
		proxyUrl.User = url.UserPassword(user, pass)
	}
	client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	return client
}
