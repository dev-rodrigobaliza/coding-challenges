package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type (
	server struct {
		Addr    string
		Healthy bool
	}

	loadBalancer struct {
		sync.RWMutex
		servers    map[int]*server
		next       int
		hcInterval time.Duration
	}
)

func New(addresses []string, healthInterval time.Duration) *loadBalancer {
	servers := make(map[int]*server)
	for i, addr := range addresses {
		server := server{
			Addr:    addr,
			Healthy: false,
		}
		servers[i] = &server
	}

	lb := &loadBalancer{
		servers:    servers,
		next:       0,
		hcInterval: healthInterval,
	}

	if len(addresses) > 0 {
		for _, server := range lb.servers {
			server.Healthy = checkHealth(server.Addr)
		}

		go lb.checkHealth()
	}

	return lb
}

func (lb *loadBalancer) Next() *server {
	var server *server
	first := lb.next

	lb.RLock()
	defer lb.RUnlock()

	for {
		svr := lb.servers[lb.next]
		lb.next++
		if lb.next >= len(lb.servers) {
			lb.next = 0
		}

		if svr.Healthy {
			server = svr
			break
		}
		if first == lb.next {
			break
		}
	}

	return server
}

func (lb *loadBalancer) checkHealth() {
	ticker := time.NewTicker(lb.hcInterval)

	for range ticker.C {
		lb.Lock()

		for _, server := range lb.servers {
			server.Healthy = checkHealth(server.Addr)
		}

		lb.Unlock()
	}
}

func checkHealth(addr string) bool {
	resp, err := http.Get(addr + "/health")
	if err != nil {
		return false
	}

	return resp.StatusCode == 200
}

var (
	lb *loadBalancer
)

func main() {
	healthInterval := flag.String("health_interval", "10s", "health check interval")

	flag.Parse()

	if *healthInterval == "" {
		fmt.Println("please specify a valid health check interval (duration)")
		os.Exit(0)
	}

	hcInterval, err := time.ParseDuration(*healthInterval)
	if err != nil {
		fmt.Println("please specify a valid health check interval (duration)")
		os.Exit(0)
	}

	lb = New([]string{"http://localhost:8001", "http://localhost:8002"}, hcInterval)

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)

	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")

	server := lb.Next()
	if server == nil {
		http.Error(w, "all servers are down now, try again later", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get(server.Addr + "/")
	if err != nil {
		msg := fmt.Sprintf("failed to get response from server: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("got bad response from server: %d", resp.StatusCode)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("could not read response body: %s", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(body))
}
