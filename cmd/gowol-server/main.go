package main

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/miguelsaliba/gowol/wol"
)

//go:embed templates/*.html
var templates embed.FS

var formTemplate = template.Must(template.ParseFS(templates, "templates/form.html"))

const (
	defaultBroadcast = "255.255.255.255"
	defaultPort      = "9"
)

var (
	errInvalidMAC       = errors.New("invalid MAC address")
	errInvalidBroadcast = errors.New("invalid broadcast address: must be an IPv4 address")
	errInvalidPort      = errors.New("invalid port: must be a number between 1 and 65535")
	errSendFailed       = errors.New("failed to send WOL packet")
)

type Data struct {
	Err       error
	Msg       string
	Mac       string
	Broadcast string
	Port      string
}

// A validated WOL destination
type target struct {
	mac  wol.MACAddress
	ip   net.IP
	port uint16
}

// Validates the form data, falling back to defaults for broadcast and port
func parseRequest(data Data) (target, error) {
	mac, err := wol.ParseMAC(data.Mac)
	if err != nil {
		return target{}, errInvalidMAC
	}

	broadcast := data.Broadcast
	if broadcast == "" {
		broadcast = defaultBroadcast
	}
	ip, err := wol.ParseIP(broadcast)
	if err != nil || ip.To4() == nil {
		return target{}, errInvalidBroadcast
	}

	port := data.Port
	if port == "" {
		port = defaultPort
	}
	parsedPort, err := strconv.ParseUint(port, 10, 16)
	if err != nil || parsedPort == 0 {
		return target{}, errInvalidPort
	}

	return target{mac: mac, ip: ip.To4(), port: uint16(parsedPort)}, nil
}

func wolHandler(w http.ResponseWriter, r *http.Request) {
	data := Data{
		Mac:       r.FormValue("mac"),
		Broadcast: r.FormValue("broadcast"),
		Port:      r.FormValue("port"),
	}

	if data.Mac == "" && data.Broadcast == "" && data.Port == "" {
		render(w, http.StatusOK, data)
		return
	}

	slog.Info("Request data", "mac", data.Mac, "broadcast", data.Broadcast, "port", data.Port)

	t, err := parseRequest(data)
	if err != nil {
		data.Err = err
		render(w, http.StatusBadRequest, data)
		return
	}

	if err := wol.Send(t.ip, t.port, t.mac); err != nil {
		slog.Error("Error sending WOL packet", "error", err)
		data.Err = errSendFailed
		render(w, http.StatusInternalServerError, data)
		return
	}

	data.Msg = "WOL successful! Bookmark this URL to use it again"
	render(w, http.StatusOK, data)
}

// Renders the form into a buffer first so a template error never produces a partial page
func render(w http.ResponseWriter, status int, data Data) {
	var buf bytes.Buffer
	if err := formTemplate.Execute(&buf, data); err != nil {
		slog.Error("Error rendering template", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

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
