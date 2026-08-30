package main

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/miguelsaliba/gowol/wol"
)

//go:embed templates/*.html
var templates embed.FS

var formTemplate = template.Must(template.ParseFS(templates, "templates/form.html"))

type Data struct {
	Err       error
	Msg       string
	Mac       string
	Broadcast string
	Port      string
}

func wolHandler(w http.ResponseWriter, r *http.Request) {
	data := Data{
		Mac:       r.FormValue("mac"),
		Broadcast: r.FormValue("broadcast"),
		Port:      r.FormValue("port"),
	}

	if data.Mac == "" && data.Broadcast == "" && data.Port == "" {
		writeHTML(w, data)
		return
	}

	slog.Info("Request data", "mac", data.Mac, "broadcast", data.Broadcast, "port", data.Port)

	parsedMac, err := wol.ParseMAC(data.Mac)
	if err != nil {
		data.Err = err
		writeHTML(w, data)
		return
	}

	parsedIP, err := wol.ParseIP(data.Broadcast)
	if err != nil {
		data.Err = err
		writeHTML(w, data)
		return
	}

	parsedPort64, err := strconv.ParseUint(data.Port, 10, 16)
	if err != nil {
		data.Err = err
		writeHTML(w, data)
		return
	}
	parsedPort := uint16(parsedPort64)

	err = wol.Send(parsedIP, parsedPort, parsedMac)
	if err != nil {
		data.Err = err
		writeHTML(w, data)
		return
	}

	data.Msg = "WOL successful! Bookmark this URL to use it again"
	writeHTML(w, data)
}

func writeHTML(w http.ResponseWriter, data Data) {
	if err := formTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Returns environment variable if exists, otherwise returns fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	port := getEnv("GOWOL_PORT", "2611")

	mux := http.NewServeMux()
	mux.HandleFunc("/", wolHandler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	slog.Info("Serving on address http://127.0.0.1:" + port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
