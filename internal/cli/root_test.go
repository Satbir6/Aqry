package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"aqry/internal/dns"
	"aqry/internal/tui"
)

type mockResolver struct {
	result  dns.LookupResult
	err     error
	request dns.LookupRequest
}

func (r *mockResolver) Lookup(_ context.Context, request dns.LookupRequest) (dns.LookupResult, error) {
	r.request = request
	return r.result, r.err
}

func executeCommand(t *testing.T, resolver dns.Resolver, runner TUIRunner, args ...string) (string, string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	command := NewRootCommand(Dependencies{
		Resolver: resolver,
		Out:      &stdout,
		Err:      &stderr,
		RunTUI:   runner,
	})
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), stderr.String(), err
}

func TestRootCommandDefaultLookup(t *testing.T) {
	t.Parallel()

	resolver := &mockResolver{result: dns.LookupResult{
		Domain: "example.com", RecordType: dns.RecordA, Resolver: "system",
		Records: []dns.Record{{Type: dns.RecordA, Value: "192.0.2.1"}},
	}}
	stdout, _, err := executeCommand(t, resolver, nil, "example.com")
	if err != nil {
		t.Fatal(err)
	}
	if stdout != "192.0.2.1\n" {
		t.Fatalf("stdout = %q", stdout)
	}
	if resolver.request.RecordType != dns.RecordA || resolver.request.All {
		t.Fatalf("request = %#v", resolver.request)
	}
}

func TestRootCommandForwardsFlags(t *testing.T) {
	t.Parallel()

	resolver := &mockResolver{result: dns.LookupResult{
		Domain: "example.com", RecordType: dns.RecordMX, Resolver: "1.1.1.1",
		Records: []dns.Record{
			{Type: dns.RecordMX, Value: "mx1.example.com", Priority: 10},
			{Type: dns.RecordMX, Value: "mx2.example.com", Priority: 20},
		},
	}}
	stdout, _, err := executeCommand(t, resolver, nil,
		"example.com", "--type", "mx", "--all", "--server", "1.1.1.1", "--timeout", "7",
	)
	if err != nil {
		t.Fatal(err)
	}
	if stdout != "10 mx1.example.com\n20 mx2.example.com\n" {
		t.Fatalf("stdout = %q", stdout)
	}
	if resolver.request.RecordType != dns.RecordMX || !resolver.request.All || resolver.request.Server != "1.1.1.1" {
		t.Fatalf("request = %#v", resolver.request)
	}
	if resolver.request.Timeout.Seconds() != 7 {
		t.Fatalf("timeout = %s", resolver.request.Timeout)
	}
}

func TestRootCommandSelectsInteractiveModes(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		args   []string
		domain string
	}{
		{name: "no arguments", args: nil},
		{name: "forced with domain", args: []string{"-i", "example.com"}, domain: "example.com"},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var received tui.Options
			runner := func(options tui.Options) error {
				received = options
				return nil
			}
			_, _, err := executeCommand(t, &mockResolver{}, runner, test.args...)
			if err != nil {
				t.Fatal(err)
			}
			if received.Domain != test.domain || received.RecordType != dns.RecordA {
				t.Fatalf("TUI options = %#v", received)
			}
		})
	}
}

func TestRootCommandVersion(t *testing.T) {
	t.Parallel()

	stdout, _, err := executeCommand(t, &mockResolver{}, nil, "--version")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(stdout, "aqry version ") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestRootCommandValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args []string
		code int
	}{
		{args: []string{"example.com", "--timeout", "0"}, code: 1},
		{args: []string{"one.example", "two.example"}, code: 1},
		{args: []string{"example.com", "--type", "SOA"}, code: 4},
		{args: []string{"example.com", "--does-not-exist"}, code: 1},
	}
	for _, test := range tests {
		_, _, err := executeCommand(t, &mockResolver{}, nil, test.args...)
		if err == nil {
			t.Fatalf("args %v returned nil error", test.args)
		}
		if actual := ExitCode(err); actual != test.code {
			t.Errorf("args %v exit code = %d, want %d", test.args, actual, test.code)
		}
	}
}

func TestExitCodeMapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err  error
		code int
	}{
		{err: nil, code: 0},
		{err: dns.NewLookupError(dns.ErrInvalidDomain, "", "", nil), code: 1},
		{err: dns.NewLookupError(dns.ErrNoRecords, "example.com", dns.RecordA, nil), code: 2},
		{err: dns.NewLookupError(dns.ErrResolverFailed, "example.com", dns.RecordA, nil), code: 2},
		{err: dns.NewLookupError(dns.ErrTimeout, "example.com", dns.RecordA, nil), code: 3},
		{err: dns.NewLookupError(dns.ErrUnsupportedType, "", dns.RecordSOA, nil), code: 4},
		{err: errors.New("unexpected"), code: 5},
	}
	for _, test := range tests {
		if actual := ExitCode(test.err); actual != test.code {
			t.Errorf("ExitCode(%v) = %d, want %d", test.err, actual, test.code)
		}
	}
}
