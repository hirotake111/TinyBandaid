package main

import (
	"fmt"
	"log"
	"net/http"

	"workspace/tinybandaid/config"
	"workspace/tinybandaid/internal/pool"
)

const CONFIG_FILE_NAME = "config.json"

// HTTP handler for testing purpose
func Pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"message": "pong"}`)
}

func main() {
	c, err := config.LoadConfigFile(CONFIG_FILE_NAME)
	if err != nil {
		panic(err)
	}

	p := pool.New(c.ServerUrls)
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", Pong)
	mux.HandleFunc("/", p.CreateHandler())
	server := http.Server{Addr: ":3000", Handler: mux}

	fmt.Println("Load balancer is up and running - http://localhost:3000")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Load balancer stopped.")
}
