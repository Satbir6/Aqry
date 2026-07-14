package dns

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

func TestSystemResolverRejectsInvalidRequestBeforeNetworkLookup(t *testing.T) {
	t.Parallel()

	resolver := NewSystemResolver()
	_, err := resolver.Lookup(context.Background(), LookupRequest{Domain: "bad domain", RecordType: RecordA})
	if !IsErrorKind(err, ErrInvalidDomain) {
		t.Fatalf("error = %v, want invalid-domain kind", err)
	}

	_, err = resolver.Lookup(context.Background(), LookupRequest{Domain: "example.com", RecordType: RecordSOA})
	if !IsErrorKind(err, ErrUnsupportedType) {
		t.Fatalf("error = %v, want unsupported-type kind", err)
	}
}

func TestValidateResolverServer(t *testing.T) {
	t.Parallel()

	for _, server := range []string{"1.1.1.1", "8.8.8.8:5353", "2001:4860:4860::8888", "[::1]:5353"} {
		if err := ValidateResolverServer(server); err != nil {
			t.Errorf("ValidateResolverServer(%q) returned %v", server, err)
		}
	}
	for _, server := range []string{"", "dns.example.com", "1.1.1.1:99999"} {
		if err := ValidateResolverServer(server); err == nil {
			t.Errorf("ValidateResolverServer(%q) returned nil", server)
		}
	}
}

func TestClassifyLookupError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	if err := classifyLookupError(ctx, "example.com", RecordA, context.DeadlineExceeded); !IsErrorKind(err, ErrTimeout) {
		t.Fatalf("deadline error classified as %v", err)
	}

	plain := errors.New("network unavailable")
	if err := classifyLookupError(context.Background(), "example.com", RecordA, plain); !IsErrorKind(err, ErrResolverFailed) {
		t.Fatalf("plain error classified as %v", err)
	}
}

func TestTrimDNSRootDotPreservesNullMXValue(t *testing.T) {
	t.Parallel()

	if actual := trimDNSRootDot("mail.example.com."); actual != "mail.example.com" {
		t.Fatalf("trimDNSRootDot returned %q", actual)
	}
	if actual := trimDNSRootDot("."); actual != "." {
		t.Fatalf("root label was changed to %q", actual)
	}
}

func TestLiveSystemResolver(t *testing.T) {
	if testing.Short() || os.Getenv("AQRY_LIVE_DNS_TEST") != "1" {
		t.Skip("set AQRY_LIVE_DNS_TEST=1 to enable the live DNS smoke test")
	}

	resolver := NewSystemResolver()
	result, err := resolver.Lookup(context.Background(), LookupRequest{
		Domain: "example.com", RecordType: RecordA, Timeout: 3 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) == 0 {
		t.Fatal("live lookup returned no records")
	}
}
