package main

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/miguelsaliba/gowol/wol"
)

//go:embed templates/*.html
var templates embed.FS

var formTemplate = template.Must(template.ParseFS(templates, "templates/form.html"))

var errSendFailed = errors.New("failed to send WOL packet")

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
