package dns

import (
	"strings"
	"testing"
)

func TestNormalizeDomain(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"https://Google.com/search?q=x": "google.com",
		"http://example.com/about":      "example.com",
		"www.example.com/":              "www.example.com",
		" GOOGLE.COM ":                  "google.com",
		"api.example.com?debug=true":    "api.example.com",
		"example.com:8080/path":         "example.com",
		"example.com.":                  "example.com",
	}
	for input, expected := range tests {
		input, expected := input, expected
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if actual := NormalizeDomain(input); actual != expected {
				t.Fatalf("NormalizeDomain(%q) = %q, want %q", input, actual, expected)
			}
		})
	}
}

func TestValidateDomain(t *testing.T) {
	t.Parallel()

	valid := []string{"google.com", "www.google.com", "api.example.co.uk", "localhost"}
	for _, domain := range valid {
		if err := ValidateDomain(domain); err != nil {
			t.Errorf("ValidateDomain(%q) returned %v", domain, err)
		}
	}

	invalid := []string{
		"",
		"has space.com",
		"bad_domain.com",
		"-start.example",
		"end-.example",
		"two..dots.example",
		"192.0.2.1",
		strings.Repeat("a", 64) + ".example",
		strings.Repeat("a.", 127) + "example",
	}
	for _, domain := range invalid {
		err := ValidateDomain(domain)
		if !IsErrorKind(err, ErrInvalidDomain) {
			t.Errorf("ValidateDomain(%q) error = %v, want invalid-domain kind", domain, err)
		}
	}
}

func TestNormalizeAndValidateDomainRejectsProtocolOnlyInput(t *testing.T) {
	t.Parallel()

	_, err := NormalizeAndValidateDomain("https://")
	if !IsErrorKind(err, ErrInvalidDomain) {
		t.Fatalf("error = %v, want invalid-domain kind", err)
	}
}
