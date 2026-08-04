package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/miguelsaliba/gowol/wol"
)

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

	slog.Debug("Request data", "mac", data.Mac, "broadcast", data.Broadcast, "port", data.Port)

	if data.Mac == "" && data.Broadcast == "" && data.Port == "" {
		writeHtml(w, data)
		return
	}

	parsedMac, err := wol.ParseMac(data.Mac)
	if err != nil {
		data.Err = err
		writeHtml(w, data)
		return
	}

	parsedIP, err := wol.ParseIP(data.Broadcast)
	if err != nil {
		data.Err = err
		writeHtml(w, data)
		return
	}

	parsedPort64, err := strconv.ParseUint(data.Port, 10, 16)
	if err != nil {
		data.Err = err
		writeHtml(w, data)
		return
	}
	parsedPort := uint16(parsedPort64)

	err = wol.Wol(parsedIP, parsedPort, parsedMac)
	if err != nil {
		data.Err = err
		writeHtml(w, data)
		return
	}

	data.Msg = "WOL successful! Bookmark this URL to use it again"
	writeHtml(w, data)
}

func writeHtml(w http.ResponseWriter, data Data) {
	t, err := template.ParseFiles("tmpl/form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, data)
}

func main() {
	port := "8080"
	http.HandleFunc("/", wolHandler)
	slog.Info("Serving on address http://127.0.0.1:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
