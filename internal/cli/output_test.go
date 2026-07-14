package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"aqry/internal/dns"
)

func TestFormatPlain(t *testing.T) {
	t.Parallel()

	result := dns.LookupResult{
		Domain: "example.com", RecordType: dns.RecordA,
		Records: []dns.Record{
			{Type: dns.RecordA, Value: "192.0.2.1"},
			{Type: dns.RecordA, Value: "192.0.2.2"},
		},
	}
	first, err := FormatPlain(result, false)
	if err != nil {
		t.Fatal(err)
	}
	if first != "192.0.2.1" {
		t.Fatalf("default output = %q", first)
	}

	all, err := FormatPlain(result, true)
	if err != nil {
		t.Fatal(err)
	}
	if all != "192.0.2.1\n192.0.2.2" {
		t.Fatalf("all output = %q", all)
	}
}

func TestFormatPlainMXIncludesPriority(t *testing.T) {
	t.Parallel()

	output, err := FormatPlain(dns.LookupResult{Records: []dns.Record{
		{Type: dns.RecordMX, Value: "mail.example.com", Priority: 10},
	}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if output != "10 mail.example.com" {
		t.Fatalf("MX output = %q", output)
	}
}

func TestWriteResultJSONUsesStructuredResult(t *testing.T) {
	t.Parallel()

	expected := dns.LookupResult{
		Domain: "example.com", RecordType: dns.RecordA, Resolver: "system",
		Records: []dns.Record{{Type: dns.RecordA, Name: "example.com", Value: "192.0.2.1"}},
	}
	var output bytes.Buffer
	if err := WriteResult(&output, expected, false, true); err != nil {
		t.Fatal(err)
	}
	var actual dns.LookupResult
	if err := json.Unmarshal(output.Bytes(), &actual); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if actual.Domain != expected.Domain || actual.RecordType != expected.RecordType || len(actual.Records) != 1 {
		t.Fatalf("decoded output = %#v", actual)
	}
}

func TestFormatPlainRejectsEmptyResult(t *testing.T) {
	t.Parallel()

	_, err := FormatPlain(dns.LookupResult{Domain: "example.com", RecordType: dns.RecordA}, false)
	if !dns.IsErrorKind(err, dns.ErrNoRecords) {
		t.Fatalf("error = %v, want no-records kind", err)
	}
}
