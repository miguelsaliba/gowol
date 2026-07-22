package main

import (
	"log/slog"
	"net"

	"github.com/miguelsaliba/gowol/wol"
)

const (
	mac       = "112233445566"
	broadcast = "255.255.255.255"
	port      = 9
)

func main() {
	ip := net.ParseIP(broadcast)
	if ip == nil {
		slog.Error("Invalid broadcast IP address", "ip", broadcast)
	}

	macAddress, err := wol.ParseMac(mac)
	if err != nil {
		slog.Error("Invalid MAC address", "mac", mac, "error", err)
	}

	err = wol.Wol(ip, port, macAddress)
	if err != nil {
		slog.Error("Failed to send WOL packet", "error", err)
	}
}
