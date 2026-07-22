package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/miguelsaliba/gowol/wol"
)

const (
	defaultBroadcast = "255.255.255.255"
	defaultPort      = 9
)

func main() {
	port := flag.Uint("port", defaultPort, "Port to send the WOL packet to")
	broadcast := flag.String("broadcast", defaultBroadcast, "Broadcast IP address to send the WOL packet to")
	help := flag.Bool("help", false, "Show usage information")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	if *verbose {
		options := &slog.HandlerOptions{Level: slog.LevelDebug}
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, options)))
	}

	if flag.NArg() < 1 || *help {
		printUsage()
		return
	}

	mac := flag.Arg(0)

	ip := net.ParseIP(*broadcast)
	if ip == nil {
		slog.Error("Invalid broadcast IP address", "ip", broadcast)
		os.Exit(1)
	}

	macAddress, err := wol.ParseMac(mac)
	if err != nil {
		slog.Error("Invalid MAC address", "mac", mac, "error", err)
		os.Exit(1)
	}

	err = wol.Wol(ip, uint16(*port), macAddress)
	if err != nil {
		slog.Error("Failed to send WOL packet", "error", err)
		os.Exit(1)
	}
}

func printUsage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "Usage: gowol <MAC_ADDRESS>\n")
	flag.PrintDefaults()
}
