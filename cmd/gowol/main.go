package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/miguelsaliba/gowol/wol"
)

const (
	defaultBroadcast = "255.255.255.255"
	defaultPort      = 9
)

type cliFlags struct {
	port      *uint
	broadcast *string
	help      *bool
	verbose   *bool
}

func newFlagSet() (*flag.FlagSet, *cliFlags) {
	fs := flag.NewFlagSet("gowol", flag.ContinueOnError)

	flags := &cliFlags{
		port:      fs.UintP("port", "p", defaultPort, "Port to send the WOL packet to"),
		broadcast: fs.StringP("broadcast", "b", defaultBroadcast, "Broadcast IP address to send the WOL packet to"),
		help:      fs.BoolP("help", "h", false, "Show usage information"),
		verbose:   fs.BoolP("verbose", "v", false, "Enable verbose logging"),
	}

	return fs, flags
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs, flags := newFlagSet()
	fs.SetOutput(stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *flags.verbose {
		options := &slog.HandlerOptions{Level: slog.LevelDebug}
		slog.SetDefault(slog.New(slog.NewTextHandler(stderr, options)))
	}

	if fs.NArg() < 1 || *flags.help {
		printUsage(fs, stdout)
		return 0
	}

	ip, err := wol.ParseIP(*flags.broadcast)
	if err != nil {
		slog.Error("Invalid broadcast address", "error", err)
		return 1
	}

	mac, err := wol.ParseMAC(fs.Arg(0))
	if err != nil {
		slog.Error("Invalid MAC address", "error", err)
		return 1
	}

	if err := wol.Send(ip, uint16(*flags.port), mac); err != nil {
		slog.Error("Failed to send WOL packet", "error", err)
		return 1
	}

	return 0
}

func printUsage(fs *flag.FlagSet, w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "\tgowol <mac_address> [arguments]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Arguments:")
	fmt.Fprint(w, fs.FlagUsages())
}
