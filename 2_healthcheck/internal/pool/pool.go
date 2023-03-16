package pool

import (
	"errors"
	"log"
	"net/http"

	"workspace/tinybandaid/internal/models"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

type HealthCheck struct {
	index   int
	isAlive bool
}

type AvailableBackend struct {
	backend *models.Backend
	err     error
}

type HealthCheckFunc func(
	getClonedBackends func() []models.Backend,
	updateBackend func(idx int, isAlive bool),
)

type LB interface {
	GetAvailableIndex(currentIndex int, backends []*models.Backend) (int, error)
}

type Pool struct {
	backends            []*models.Backend     // array of backends
	currentIndex        int                   // current index of serving backend
	reqStream           chan interface{}      // used to request next available backend
	resStream           chan AvailableBackend // used to share next available backend
	hcStream            chan HealthCheck      // health check stream
	cloneRequestStream  chan interface{}      // use to request an array of cloned backends
	cloneResponseStream chan []models.Backend // used to send an array of cloned backends
	lb                  LB
}

// Returns a pointer to next available backend server
func (p *Pool) nextBackend() (*models.Backend, error) {
	index, err := p.lb.GetAvailableIndex(p.currentIndex, p.backends)
	if err != nil {
		return &models.Backend{}, errors.New("all backend servers are down")
	}
	p.currentIndex = index
	return p.backends[p.currentIndex], nil
	// // Round robin loadbalancing
	// for i := 0; i < len(p.backends); i++ {
	// p.currentIndex = (p.currentIndex + 1) % len(p.backends)
	// if p.backends[p.currentIndex].IsAlive {
	// return p.backends[p.currentIndex], nil
	// }
	// }
	// log.Println("All backend servers are down!")
	// return &models.Backend{}, errors.New("all backend servers are down")
}

// Returns a pointer for a new server pool
func New(serverUrls []string, loadbalancingMethod LB) *Pool {
	reqStream := make(chan interface{})
	resStream := make(chan AvailableBackend)
	hcStream := make(chan HealthCheck)
	cloneRequestStream := make(chan interface{})
	cloneResponseStream := make(chan []models.Backend)
	backends := make([]*models.Backend, len(serverUrls))
	for i, serverUrl := range serverUrls {
		backends[i] = models.NewBackend(serverUrl)
		log.Printf("server #%d: %s registered.\n", i, backends[i].Url)
	}
	pool := &Pool{
		backends:            backends,
		currentIndex:        len(backends) - 1,
		reqStream:           reqStream,
		resStream:           resStream,
		hcStream:            hcStream,
		cloneRequestStream:  cloneRequestStream,
		cloneResponseStream: cloneResponseStream,
		lb:                  loadbalancingMethod,
	}

	go pool.poolManager()
	return pool
}

// Returns an HTTP handler for load balancing
func (p *Pool) CreateHandler() HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get next available backend
		p.reqStream <- struct{}{}
		available := <-p.resStream
		if available.err != nil {
			http.Error(w, "BadGateway!!", http.StatusBadGateway)
			return
		}
		available.backend.ReverseProxy.ServeHTTP(w, r)
	}
}

// Updates status of backends[i] only if changed
func (p *Pool) UpdateBackendHealth(i int, isAlive bool) {
	p.hcStream <- HealthCheck{index: i, isAlive: isAlive}
}

// Returns an array of cloned backends
func (p *Pool) GetClonedBackends() []models.Backend {
	// Send a request
	p.cloneRequestStream <- struct{}{}
	// Receive a response
	return <-p.cloneResponseStream
}

// poolManager responsible for managing backends and other resources in Pool
func (p *Pool) poolManager() {
	for {
		select {
		// Health check
		case hc := <-p.hcStream:
			p.backends[hc.index].IsAlive = hc.isAlive
		// Available backend
		case <-p.reqStream:
			backend, err := p.nextBackend()
			p.resStream <- AvailableBackend{backend: backend, err: err}
		// Clone backends
		case <-p.cloneRequestStream:
			cloned := make([]models.Backend, len(p.backends))
			for i, b := range p.backends {
				cloned[i] = *b
			}
			p.cloneResponseStream <- cloned
		}
	}
}
