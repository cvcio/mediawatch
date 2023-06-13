package proxy

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// CreateProxy returns an http.Client with a proxy
func CreateProxy(proxylist []string, user, pass string) *http.Client {
	client := &http.Client{Timeout: 30 * time.Second}
	proxyUrl := &url.URL{Scheme: "http", Host: proxylist[rand.Intn(len(proxylist))]}
	if user != "" && pass != "" {
		proxyUrl.User = url.UserPassword(user, pass)
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyUrl),
	}
	client.Transport = transport
	return client
}

// CreateClient returns an http.Client
func CreateClient() *http.Client {
	client := &http.Client{Timeout: 30 * time.Second}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.Transport = transport
	return client
}
