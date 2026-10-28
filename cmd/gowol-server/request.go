package main

import (
	"errors"
	"net"
	"strconv"

	"github.com/miguelsaliba/gowol/wol"
)

const (
	defaultBroadcast = "255.255.255.255"
	defaultPort      = "9"
)

var (
	errInvalidMAC       = errors.New("invalid MAC address")
	errInvalidBroadcast = errors.New("invalid broadcast address: must be an IPv4 address")
	errInvalidPort      = errors.New("invalid port: must be a number between 1 and 65535")
)

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
