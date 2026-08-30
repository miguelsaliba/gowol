package wol

import (
	"bytes"
	"errors"
	"net"
	"testing"
	"time"
)

func TestParseMac(t *testing.T) {
	cases := []struct {
		Mac    string
		Parsed MACAddress
		Error  error
	}{
		{Mac: "11:22:33:44:55:66", Parsed: MACAddress{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}, Error: nil},
		{Mac: "FF:FF:FF:FF:FF:FF", Parsed: MACAddress{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, Error: nil},
		{Mac: "FF-FF-FF-FF-FF-FF", Parsed: MACAddress{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, Error: nil},
		{Mac: "111111111111", Parsed: MACAddress{0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, Error: nil},
		{Mac: "11:22:33:44:55:66:77:88", Parsed: MACAddress{}, Error: ErrInvalidMAC},
		{Mac: "11:22", Parsed: MACAddress{}, Error: ErrInvalidMAC},
		{Mac: "11:22:33:44", Parsed: MACAddress{}, Error: ErrInvalidMAC},
	}
	for _, tt := range cases {
		t.Run(tt.Mac, func(t *testing.T) {
			got, err := ParseMAC(tt.Mac)

			switch {
			case tt.Error == nil && err != nil:
				t.Errorf("unexpected error %v", err)

			case tt.Error != nil && err == nil:
				t.Errorf("expected error %v but got nil", tt.Error)

			case !errors.Is(err, tt.Error):
				t.Errorf("expected error %v but got %v", tt.Error, err)

			case tt.Parsed != got:
				t.Errorf("expected %v but got %v", tt.Parsed, got)
			}
		})

	}
}

func TestParseIP(t *testing.T) {
	cases := []struct {
		input   string
		want    net.IP
		wantErr bool
	}{
		{input: "192.168.1.1", want: net.ParseIP("192.168.1.1")},
		{input: "255.255.255.255", want: net.ParseIP("255.255.255.255")},
		{input: "2001:db8::1", want: net.ParseIP("2001:db8::1")},
		{input: "::1", want: net.ParseIP("::1")},
		{input: "", wantErr: true},
		{input: "not-an-ip", wantErr: true},
		{input: "256.256.256.256", wantErr: true},
		{input: "192.168.1", wantErr: true},
		{input: "192.168.1.1.1", wantErr: true},
		{input: "  192.168.1.1  ", wantErr: true},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseIP(tt.input)

			switch {
			case tt.wantErr && err == nil:
				t.Fatalf("ParseIP(%q) expected error, got nil", tt.input)

			case tt.wantErr && !errors.Is(err, ErrInvalidIP):
				t.Errorf("ParseIP(%q) error = %v, want wrapped ErrInvalidIP", tt.input, err)

			case !tt.wantErr && err != nil:
				t.Fatalf("ParseIP(%q) unexpected error: %v", tt.input, err)

			case !tt.wantErr && !got.Equal(tt.want):
				t.Errorf("ParseIP(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCreateMagicPacket(t *testing.T) {
	mac := MACAddress{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x01}

	packet := createMagicPacket(mac)

	expected := append(
		bytes.Repeat([]byte{0xFF}, 6),
		bytes.Repeat(mac[:], 16)...,
	)

	if !bytes.Equal(packet, expected) {
		t.Fatalf("packet mismatch:\ngot:  % X\nwant: % X", packet, expected)
	}

	if len(packet) != 102 {
		t.Fatalf("expected packet length 102, got %d", len(packet))
	}
}

func TestWol_SendsCorrectPacket(t *testing.T) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
	})
	if err != nil {
		t.Fatalf("failed to start UDP listener: %v", err)
	}
	defer listener.Close()

	addr := listener.LocalAddr().(*net.UDPAddr)
	mac := MACAddress{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x01}

	if err := Send(net.ParseIP("127.0.0.1"), uint16(addr.Port), mac); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	buf := make([]byte, 128)
	listener.SetReadDeadline(time.Now().Add(2 * time.Second))

	n, _, err := listener.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("failed to read UDP packet: %v", err)
	}

	expected := createMagicPacket(mac)

	if !bytes.Equal(buf[:n], expected) {
		t.Fatalf("packet mismatch:\ngot:  % X\nwant: % X", buf[:n], expected)
	}
}
