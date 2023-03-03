package pool

import (
	"log"
	"net/http"

	"github.com/hirotake111/go-toy-lb/internal/models"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

type Pool struct {
	backends     []*models.Backend // array of backends
	currentIndex uint64            // current index of serving backend
}

// Returns a pointer for next available backend server
func (p *Pool) NextBackend() *models.Backend {
	defer func() {
		p.currentIndex = (p.currentIndex + 1) % uint64(len(p.backends))
	}()
	return p.backends[p.currentIndex]
}

// Returns a pointer for a new server pool
func New(serverUrls []string) *Pool {
	backends := make([]*models.Backend, len(serverUrls))
	p := Pool{backends: backends, currentIndex: 0}
	for i, serverUrl := range serverUrls {
		p.backends[i] = models.NewBackend(serverUrl)
		log.Printf("server #%d: %s registered.\n", i, serverUrl)
	}
	return &p
}

// Return an HTTP handler for load balancing
func (p *Pool) CreateHandler() HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		b := p.NextBackend()
		b.ReverseProxy.ServeHTTP(w, r)
	}
}
