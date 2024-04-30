package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	parser "github.com/openvenues/gopostal/parser"
)

func main() {
	slog.Info("Starting server", "port", 9876)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthy", healthy)
	mux.HandleFunc("/parse/{address...}", parse)

	http.ListenAndServe(":9876", mux)
}

func healthy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I'm healthy"))
}

func parse(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.Method != "GET" {
		slog.Error("Method not allowed parsing address", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.PathValue("address")

	parsed := parser.ParseAddress(address)

	data, err := json.Marshal(parsed)
	if err != nil {
		slog.Error("Error parsing address", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	slog.Info("Address parsed", "address", address, "parsed", parsed, "elapsed", time.Since(start))
}
