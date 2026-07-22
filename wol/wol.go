package wol

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net"
)

type MACAddress [6]byte

var ErrInvalidMAC = errors.New("invalid Wake-on-LAN MAC address")
var ErrInvalidIP = errors.New("invalid IP address")

func ParseMac(mac string) (MACAddress, error) {
	parsed, err := net.ParseMAC(mac)
	if err != nil {
		return MACAddress{}, fmt.Errorf("%w: %v", ErrInvalidMAC, err)
	}

	if len(parsed) != 6 {
		return MACAddress{}, fmt.Errorf("%w: %v", ErrInvalidMAC, "MAC address must be 6 bytes long")
	}

	return MACAddress(parsed), nil
}

func Wol(ip net.IP, port uint16, mac MACAddress) error {
	if ip == nil {
		return ErrInvalidIP
	}

	udpAddr := &net.UDPAddr{
		IP:   ip,
		Port: int(port),
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %w", err)
	}
	defer conn.Close()

	syncStream := bytes.Repeat([]byte{0xFF}, 6)
	packet := append(syncStream, bytes.Repeat(mac[:], 16)...)

	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("failed to send WOL packet: %w", err)
	}

	slog.Debug("WOL packet successfully sent to", "IP", ip, "Port", port, "MAC", mac)

	return nil
}

func (m MACAddress) String() string {
	return net.HardwareAddr(m[:]).String()
}
