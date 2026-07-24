package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Version is set at build time via ldflags
var Version = "dev"

type InfoResponse struct {
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Host      string `json:"host"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	resp := InfoResponse{
		Service:   "cicd-demo",
		Version:   Version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Host:      host,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{Status: "healthy"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", infoHandler)
	mux.HandleFunc("/healthz", healthHandler)
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := setupRoutes()

	log.Printf("Starting cicd-demo %s on :%s", Version, port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
