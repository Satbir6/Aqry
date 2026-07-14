package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"aqry/internal/dns"
)

func WriteResult(writer io.Writer, result dns.LookupResult, all, asJSON bool) error {
	if asJSON {
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	output, err := FormatPlain(result, all)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, output)
	return err
}

func FormatPlain(result dns.LookupResult, all bool) (string, error) {
	if len(result.Records) == 0 {
		return "", dns.NewLookupError(dns.ErrNoRecords, result.Domain, result.RecordType, nil)
	}

	limit := 1
	if all {
		limit = len(result.Records)
	}
	lines := make([]string, 0, limit)
	for _, record := range result.Records[:limit] {
		if record.Type == dns.RecordMX {
			lines = append(lines, fmt.Sprintf("%d %s", record.Priority, record.Value))
		} else {
			lines = append(lines, record.Value)
		}
	}

	return strings.Join(lines, "\n"), nil
}
