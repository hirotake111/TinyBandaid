package pool

import (
	"log"
	"net/http"

	"workspace/tinybandaid/internal/models"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

type Pool struct {
	backends     []*models.Backend // array of backends
	currentIndex uint64            // current index of serving backend
}

// Returns a pointer to next available backend server
func (p *Pool) nextBackend() *models.Backend {
	p.currentIndex = (p.currentIndex + 1) % uint64(len(p.backends))
	return p.backends[p.currentIndex]
}

// Returns a pointer for a new server pool
func New(serverUrls []string) *Pool {
	backends := make([]*models.Backend, len(serverUrls))
	for i, serverUrl := range serverUrls {
		backends[i] = models.NewBackend(serverUrl)
		log.Printf("server #%d: %s registered.\n", i, backends[i].Url)
	}
	return &Pool{
		backends:     backends,
		currentIndex: uint64(len(backends) - 1),
	}
}

// Returns an HTTP handler for load balancing
func (p *Pool) CreateHandler() HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		b := p.nextBackend()
		b.ReverseProxy.ServeHTTP(w, r)
	}
}
