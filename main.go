package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/miguelsaliba/gowol/wol"
)

const (
	defaultBroadcast = "255.255.255.255"
	defaultPort      = 9
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("gowol", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// flags
	port := fs.UintP("port", "p", defaultPort, "Port to send the WOL packet to")
	broadcast := fs.StringP("broadcast", "b", defaultBroadcast, "Broadcast IP address to send the WOL packet to")
	help := fs.BoolP("help", "h", false, "Show usage information")
	verbose := fs.BoolP("verbose", "v", false, "Enable verbose logging")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *verbose {
		options := &slog.HandlerOptions{Level: slog.LevelDebug}
		slog.SetDefault(slog.New(slog.NewTextHandler(stderr, options)))
	}

	if fs.NArg() < 1 || *help {
		printUsage(fs, stdout)
		return 0
	}

	macStr := fs.Arg(0)

	ip, mac, err := validateInputs(*broadcast, macStr)
	if err != nil {
		slog.Error(err.Error())
		return 1
	}

	if err := wol.Wol(ip, uint16(*port), mac); err != nil {
		slog.Error("Failed to send WOL packet", "error", err)
		return 1
	}

	return 0
}

func validateInputs(broadcastStr, macStr string) (net.IP, wol.MACAddress, error) {
	ip := net.ParseIP(broadcastStr)
	if ip == nil {
		return nil, wol.MACAddress{}, fmt.Errorf("invalid broadcast IP address: %s", broadcastStr)
	}

	mac, err := wol.ParseMac(macStr)
	if err != nil {
		return nil, wol.MACAddress{}, fmt.Errorf("invalid MAC address %q: %w", macStr, err)
	}

	return ip, mac, nil
}

func printUsage(fs *flag.FlagSet, w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "\tgowol <mac_address> [arguments]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Arguments:")
	fmt.Fprint(w, fs.FlagUsages())
}
