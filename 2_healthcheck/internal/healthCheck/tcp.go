package healthCheck

import (
	"log"
	"time"
	"workspace/tinybandaid/internal/models"
)

const DURATION = 3 * time.Second

type pool interface {
	GetClonedBackends() []models.Backend
	UpdateBackendHealth(idx int, isAlive bool)
}

func TcpHealthCheck(p pool) {
	ticker := time.NewTicker(DURATION)
	log.Println("Health check started")
	for range ticker.C {
		clonedBackends := p.GetClonedBackends()
		for i, b := range clonedBackends {
			p.UpdateBackendHealth(i, isBackendAlive(b.Url))
		}
	}
}
