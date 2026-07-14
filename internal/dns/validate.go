package dns

import (
	"fmt"
	"net"
	"strings"
	"unicode"
)

func ValidateDomain(domain string) error {
	if domain == "" {
		return NewLookupError(ErrInvalidDomain, domain, "", fmt.Errorf("enter a domain to continue"))
	}
	if len(domain) > 253 {
		return NewLookupError(ErrInvalidDomain, domain, "", fmt.Errorf("domain exceeds 253 characters"))
	}
	if net.ParseIP(domain) != nil {
		return NewLookupError(ErrInvalidDomain, domain, "", fmt.Errorf("an IP address is not a domain name"))
	}
	if strings.IndexFunc(domain, unicode.IsSpace) >= 0 {
		return NewLookupError(ErrInvalidDomain, domain, "", fmt.Errorf("domain cannot contain spaces"))
	}

	for _, label := range strings.Split(domain, ".") {
		if label == "" {
			return NewLookupError(ErrInvalidDomain, domain, "", fmt.Errorf("domain contains an empty label"))
		}
		if len(label) > 63 {
			return NewLookupError(ErrInvalidDomain, domain, "", fmt.Errorf("domain label exceeds 63 characters"))
		}
		if label[0] == '-' || label[len(label)-1] == '-' {
			return NewLookupError(ErrInvalidDomain, domain, "", fmt.Errorf("domain labels cannot start or end with a hyphen"))
		}
		for _, character := range label {
			if (character >= 'a' && character <= 'z') ||
				(character >= '0' && character <= '9') || character == '-' {
				continue
			}
			return NewLookupError(ErrInvalidDomain, domain, "", fmt.Errorf("domain contains invalid character %q", character))
		}
	}

	return nil
}
