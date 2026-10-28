package main

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func doRequest(t *testing.T, form url.Values) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/?"+form.Encode(), nil)
	rec := httptest.NewRecorder()
	wolHandler(rec, req)
	return rec
}

func TestWolHandler(t *testing.T) {
	tests := []struct {
		name       string
		form       url.Values
		wantStatus int
		wantBody   string
	}{
		{
			name:       "empty form shows blank page",
			form:       url.Values{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid mac address",
			form:       url.Values{"mac": {"not-a-mac"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   errInvalidMAC.Error(),
		},
		{
			name:       "invalid broadcast address",
			form:       url.Values{"mac": {"00:11:22:33:44:55"}, "broadcast": {"not-an-ip"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   errInvalidBroadcast.Error(),
		},
		{
			name:       "invalid port",
			form:       url.Values{"mac": {"00:11:22:33:44:55"}, "port": {"0"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   errInvalidPort.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(t, tt.form)

			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("expected body to contain %q, got %q", tt.wantBody, rec.Body.String())
			}
		})
	}
}

func TestWolHandlerSendsMagicPacket(t *testing.T) {
	conn, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	port := conn.LocalAddr().(*net.UDPAddr).Port

	rec := doRequest(t, url.Values{
		"mac":       {"00:11:22:33:44:55"},
		"broadcast": {"127.0.0.1"},
		"port":      {strconv.Itoa(port)},
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "WOL successful!") {
		t.Errorf("expected success message, got %q", rec.Body.String())
	}

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(time.Second))
	n, _, err := conn.ReadFrom(buf)
	if err != nil {
		t.Fatalf("expected a packet, got error: %v", err)
	}

	mac := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	want := append(bytes.Repeat([]byte{0xFF}, 6), bytes.Repeat(mac, 16)...)
	if !bytes.Equal(buf[:n], want) {
		t.Errorf("expected magic packet %x, got %x", want, buf[:n])
	}
}
