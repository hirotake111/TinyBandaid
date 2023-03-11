package models

import (
	"fmt"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	Url          *url.URL               // Backend's URL
	ReverseProxy *httputil.ReverseProxy // reverse proxy instance
	IsAlive      bool                   // false if health check failed
}

// Returns a pointer to a new backend server
func NewBackend(serverUrl string) *Backend {
	url, err := url.Parse(serverUrl)
	if err != nil {
		panic(fmt.Sprintf("Cannot instantiate a server with URL %s", serverUrl))
	}
	return &Backend{
		Url:          url,
		ReverseProxy: httputil.NewSingleHostReverseProxy(url),
		IsAlive:      false,
	}
}
