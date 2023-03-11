package healthCheck

import (
	"net"
	"net/url"
	"time"
)

const HEALTH_CHECK_TIMEOUT = 2 * time.Second

func isBackendAlive(u *url.URL) bool {
	// Assign default port 80 if not assigned
	port := u.Port()
	if port == "" {
		port = "80"
	}
	address := u.Host + ":" + port
	con, err := net.DialTimeout("tcp", address, HEALTH_CHECK_TIMEOUT)
	if err != nil {
		return false
	}
	con.Close()
	return true
}
