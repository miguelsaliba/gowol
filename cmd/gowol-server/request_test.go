package main

import (
	"errors"
	"net"
	"testing"

	"github.com/miguelsaliba/gowol/wol"
)

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
