package main

import (
	"bytes"
	"strings"
	"testing"
)

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

func TestNewFlagSet(t *testing.T) {
	var out bytes.Buffer
	fs, flags := newFlagSet()
	fs.SetOutput(&out)

	if fs.Name() != "gowol" {
		t.Errorf("expected flag set name %q, got %q", "gowol", fs.Name())
	}

	if err := fs.Parse([]string{"--port", "1234", "--broadcast", "10.0.0.255", "--help", "--verbose"}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if *flags.port != 1234 {
		t.Errorf("expected port 1234, got %d", *flags.port)
	}
	if *flags.broadcast != "10.0.0.255" {
		t.Errorf("expected broadcast 10.0.0.255, got %s", *flags.broadcast)
	}
	if !*flags.help {
		t.Error("expected help to be true")
	}
	if !*flags.verbose {
		t.Error("expected verbose to be true")
	}
}

func TestPrintUsage(t *testing.T) {
	var buf bytes.Buffer
	fs, _ := newFlagSet()
	fs.SetOutput(&buf)

	var usageBuf bytes.Buffer
	printUsage(fs, &usageBuf)

	out := usageBuf.String()
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected output to contain 'Usage:', got %q", out)
	}
	if !strings.Contains(out, "gowol <mac_address>") {
		t.Errorf("expected output to contain command usage, got %q", out)
	}
}
