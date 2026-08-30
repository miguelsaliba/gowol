package wol

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net"
)

type MACAddress [6]byte

var ErrInvalidMAC = errors.New("invalid MAC address")
var ErrInvalidIP = errors.New("invalid IP address")

func ParseMAC(mac string) (MACAddress, error) {
	parsed, err := net.ParseMAC(mac)
	if err != nil {
		return MACAddress{}, ErrInvalidMAC
	}

	if len(parsed) != 6 {
		return MACAddress{}, fmt.Errorf("%w: %v", ErrInvalidMAC, "MAC address must be 6 bytes long")
	}

	return MACAddress(parsed), nil
}

func ParseIP(ip string) (net.IP, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidIP, ip)
	}
	return parsed, nil
}

func createMagicPacket(mac MACAddress) []byte {
	packet := make([]byte, 0, 6+len(mac)*16)
	packet = append(packet, bytes.Repeat([]byte{0xFF}, 6)...)
	packet = append(packet, bytes.Repeat(mac[:], 16)...)
	return packet
}

func Send(ip net.IP, port uint16, mac MACAddress) error {
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

	packet := createMagicPacket(mac)

	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("failed to send WOL packet: %w", err)
	}

	slog.Debug("WOL packet successfully sent to", "IP", ip, "Port", port, "MAC", mac)

	return nil
}

func (m MACAddress) String() string {
	return net.HardwareAddr(m[:]).String()
}
