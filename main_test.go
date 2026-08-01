package main

import (
	"bytes"
	"strings"
	"testing"

	flag "github.com/spf13/pflag"

	"github.com/miguelsaliba/gowol/wol"
)

func TestValidateInputs(t *testing.T) {
	tests := []struct {
		name       string
		broadcast  string
		mac        string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:      "valid inputs",
			broadcast: "255.255.255.255",
			mac:       "00:11:22:33:44:55",
			wantErr:   false,
		},
		{
			name:       "invalid broadcast address",
			broadcast:  "not-an-ip",
			mac:        "00:11:22:33:44:55",
			wantErr:    true,
			wantErrMsg: "invalid broadcast IP address",
		},
		{
			name:       "invalid mac address",
			broadcast:  "255.255.255.255",
			mac:        "not-a-mac",
			wantErr:    true,
			wantErrMsg: "invalid MAC address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, mac, err := validateInputs(tt.broadcast, tt.mac)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.wantErrMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ip == nil {
				t.Error("expected non-nil IP")
			}
			if (mac == wol.MACAddress{}) {
				t.Error("expected non-nil MAC address")
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
	}{
		{
			name:       "help flag shows usage",
			args:       []string{"--help"},
			wantCode:   0,
			wantStdout: "Usage:",
		},
		{
			name:       "no args shows usage",
			args:       []string{},
			wantCode:   0,
			wantStdout: "Usage:",
		},
		{
			name:     "invalid broadcast address",
			args:     []string{"--broadcast", "not-an-ip", "00:11:22:33:44:55"},
			wantCode: 1,
		},
		{
			name:     "invalid mac address",
			args:     []string{"not-a-mac"},
			wantCode: 1,
		},
		{
			name:     "successful send",
			args:     []string{"--broadcast", "127.0.0.1", "--port", "39999", "00:11:22:33:44:55"},
			wantCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			code := run(tt.args, &stdout, &stderr)

			if code != tt.wantCode {
				t.Errorf("expected exit code %d, got %d (stderr: %s)", tt.wantCode, code, stderr.String())
			}
			if tt.wantStdout != "" && !strings.Contains(stdout.String(), tt.wantStdout) {
				t.Errorf("expected stdout to contain %q, got %q", tt.wantStdout, stdout.String())
			}
		})
	}
}

func TestPrintUsage(t *testing.T) {
	fs := newTestFlagSet()

	var buf bytes.Buffer
	printUsage(fs, &buf)

	out := buf.String()
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected output to contain 'Usage:', got %q", out)
	}
	if !strings.Contains(out, "gowol <mac_address>") {
		t.Errorf("expected output to contain command usage, got %q", out)
	}
}

// newTestFlagSet mirrors the flags declared in run, for use in
// printUsage tests without invoking run itself.
func newTestFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("gowol", flag.ContinueOnError)
	fs.UintP("port", "p", defaultPort, "Port to send the WOL packet to")
	fs.StringP("broadcast", "b", defaultBroadcast, "Broadcast IP address to send the WOL packet to")
	fs.BoolP("help", "h", false, "Show usage information")
	fs.BoolP("verbose", "v", false, "Enable verbose logging")
	return fs
}
