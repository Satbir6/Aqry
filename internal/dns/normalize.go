package dns

import (
	"net/url"
	"strings"
)

func NormalizeDomain(input string) string {
	value := strings.ToLower(strings.TrimSpace(input))
	if value == "" {
		return ""
	}

	candidate := value
	if !strings.Contains(candidate, "://") && !strings.HasPrefix(candidate, "//") {
		candidate = "//" + candidate
	}

	if parsed, err := url.Parse(candidate); err == nil {
		if host := parsed.Hostname(); host != "" {
			return strings.TrimSuffix(host, ".")
		}
	}

	if cut := strings.IndexAny(value, "/?#"); cut >= 0 {
		value = value[:cut]
	}

	return strings.TrimSuffix(value, ".")
}

func NormalizeAndValidateDomain(input string) (string, error) {
	domain := NormalizeDomain(input)
	if err := ValidateDomain(domain); err != nil {
		return "", err
	}

	return domain, nil
}
