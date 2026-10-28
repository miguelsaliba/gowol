package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Returns environment variable if exists, otherwise returns fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	var level slog.Level
	if err := level.UnmarshalText([]byte(getEnv("GOWOL_LOG_LEVEL", "info"))); err != nil {
		slog.Error("Invalid GOWOL_LOG_LEVEL", "error", err)
		os.Exit(1)
	}
	slog.SetLogLoggerLevel(level)

	port := getEnv("GOWOL_PORT", "2611")

	http.HandleFunc("GET /{$}", wolHandler)

	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	slog.Info("Serving on all interfaces", "port", port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
