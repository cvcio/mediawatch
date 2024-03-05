package proxy

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// CreateProxyFromList returns an http.Client with a proxy
func CreateProxyFromList(proxylist []string, user, pass string) *http.Client {
	proxyUrl := &url.URL{Scheme: "http", Host: proxylist[rand.Intn(len(proxylist))]}
	if user != "" && pass != "" {
		proxyUrl.User = url.UserPassword(user, pass)
	}
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyUrl),
		},
	}
}

// CreateProxy returns an http.Client with a proxy
func CreateProxy(u string) *http.Client {
	proxyUrl, err := url.Parse(u)
	if err != nil {
		return nil
	}
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyUrl),
		},
	}
}

// CreateClient returns an http.Client
func CreateClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
