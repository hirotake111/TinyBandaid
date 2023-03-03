package models

import (
	"fmt"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	Url          string                 // URL for server
	ReverseProxy *httputil.ReverseProxy // reverse proxy instance
}

// Returns a pointer for a new backend server
func NewBackend(serverUrl string) *Backend {
	url, err := url.Parse(serverUrl)
	if err != nil {
		panic(fmt.Sprintf("Cannot instantiate a server with URL %s", serverUrl))
	}
	return &Backend{Url: serverUrl, ReverseProxy: httputil.NewSingleHostReverseProxy(url)}
}
