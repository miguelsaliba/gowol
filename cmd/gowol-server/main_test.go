package main

import (
	"bytes"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/miguelsaliba/gowol/wol"
)

func doRequest(t *testing.T, form url.Values) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/?"+form.Encode(), nil)
	rec := httptest.NewRecorder()
	wolHandler(rec, req)
	return rec
}

func TestParseRequest(t *testing.T) {
	mac := wol.MACAddress{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}

	tests := []struct {
		name    string
		data    Data
		want    target
		wantErr error
	}{
		{
			name: "valid input",
			data: Data{Mac: "00:11:22:33:44:55", Broadcast: "192.168.1.255", Port: "7"},
			want: target{mac: mac, ip: net.IPv4(192, 168, 1, 255).To4(), port: 7},
		},
		{
			name: "defaults for broadcast and port",
			data: Data{Mac: "00:11:22:33:44:55"},
			want: target{mac: mac, ip: net.IPv4bcast.To4(), port: 9},
		},
		{name: "missing mac", data: Data{}, wantErr: errInvalidMAC},
		{name: "invalid mac", data: Data{Mac: "not-a-mac"}, wantErr: errInvalidMAC},
		{name: "invalid broadcast", data: Data{Mac: "00:11:22:33:44:55", Broadcast: "not-an-ip"}, wantErr: errInvalidBroadcast},
		{name: "ipv6 broadcast", data: Data{Mac: "00:11:22:33:44:55", Broadcast: "ff02::1"}, wantErr: errInvalidBroadcast},
		{name: "non-numeric port", data: Data{Mac: "00:11:22:33:44:55", Port: "abc"}, wantErr: errInvalidPort},
		{name: "port 0", data: Data{Mac: "00:11:22:33:44:55", Port: "0"}, wantErr: errInvalidPort},
		{name: "port 65536", data: Data{Mac: "00:11:22:33:44:55", Port: "65536"}, wantErr: errInvalidPort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRequest(tt.data)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if got.mac != tt.want.mac || !got.ip.Equal(tt.want.ip) || got.port != tt.want.port {
				t.Errorf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
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

func TestGetEnv(t *testing.T) {
	const key = "GOWOL_TEST_ENV_VAR"

	t.Run("returns fallback when unset", func(t *testing.T) {
		t.Setenv(key, "") // restores the original value after the test
		os.Unsetenv(key)
		if got := getEnv(key, "fallback"); got != "fallback" {
			t.Errorf("expected %q, got %q", "fallback", got)
		}
	})

	t.Run("returns value when set", func(t *testing.T) {
		t.Setenv(key, "value")
		if got := getEnv(key, "fallback"); got != "value" {
			t.Errorf("expected %q, got %q", "value", got)
		}
	})
}
