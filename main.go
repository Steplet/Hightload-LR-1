package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "healthcheck" {
		os.Exit(checkHealth())
	}

	podID := strings.TrimSpace(os.Getenv("POD_ID"))
	datacenter := strings.TrimSpace(os.Getenv("DATACENTER"))
	if podID == "" || datacenter == "" {
		log.Fatal("POD_ID and DATACENTER environment variables are required")
	}

	routes := itemRoutes(NewStore())
	metrics := NewMetrics(podID, datacenter, routes)
	server := newRouter(routes, metrics)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}

func checkHealth() int {
	client := http.Client{Timeout: 2 * time.Second}
	response, err := client.Get("http://127.0.0.1:8080/health")
	if err != nil {
		return 1
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
