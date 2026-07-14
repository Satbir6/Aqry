package dns

import "strings"

type RecordType string

const (
	RecordA     RecordType = "A"
	RecordAAAA  RecordType = "AAAA"
	RecordCNAME RecordType = "CNAME"
	RecordMX    RecordType = "MX"
	RecordTXT   RecordType = "TXT"
	RecordNS    RecordType = "NS"
	RecordSOA   RecordType = "SOA"
	RecordCAA   RecordType = "CAA"
)

type Record struct {
	Type     RecordType `json:"type"`
	Name     string     `json:"name"`
	Value    string     `json:"value"`
	Priority int        `json:"priority,omitempty"`
	TTL      int        `json:"ttl,omitempty"`
}

type LookupResult struct {
	Domain     string     `json:"domain"`
	RecordType RecordType `json:"type"`
	Resolver   string     `json:"resolver"`
	Records    []Record   `json:"records"`
}

func SupportedRecordTypes() []RecordType {
	return []RecordType{RecordA, RecordAAAA, RecordCNAME, RecordMX, RecordTXT, RecordNS}
}

func ParseRecordType(value string) (RecordType, error) {
	if strings.TrimSpace(value) == "" {
		return RecordA, nil
	}

	recordType := RecordType(strings.ToUpper(strings.TrimSpace(value)))
	if !recordType.Supported() {
		return "", NewLookupError(ErrUnsupportedType, "", recordType, nil)
	}

	return recordType, nil
}

func (r RecordType) Supported() bool {
	for _, supported := range SupportedRecordTypes() {
		if r == supported {
			return true
		}
	}

	return false
}

func (r RecordType) Description() string {
	switch r {
	case RecordA:
		return "IPv4 address"
	case RecordAAAA:
		return "IPv6 address"
	case RecordCNAME:
		return "Canonical name"
	case RecordMX:
		return "Mail exchange"
	case RecordTXT:
		return "Text record"
	case RecordNS:
		return "Name server"
	default:
		return "Unsupported record"
	}
}
