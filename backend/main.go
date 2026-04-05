package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"app":    "PackSmart",
		"phase":  "Module 4 Infrastructure Test",
	})
}

func main() {
	http.HandleFunc("/api/health", healthHandler)

	log.Println("PackSmart backend running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
