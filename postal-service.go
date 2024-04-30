package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	parser "github.com/openvenues/gopostal/parser"
)

func main() {
	slog.Info("Starting server", "port", 9876)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthy", healthy)
	mux.HandleFunc("/parse/{address...}", parse)

	server := &http.Server{Addr: ":9876", Handler: mux}

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		slog.Info("Stopping server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown error", "error", err)
		}
	}()

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}

func healthy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I'm healthy"))
}

func parse(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.Method != "GET" {
		slog.Error("Method not allowed parsing address", "method", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	address := r.PathValue("address")

	parsed := parser.ParseAddress(address)

	data, err := json.Marshal(parsed)
	if err != nil {
		slog.Error("Error marshalling address", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	slog.Info("Address parsed", "address", address, "parsed", parsed, "elapsed", time.Since(start))
}
