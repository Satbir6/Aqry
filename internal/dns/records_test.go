package dns

import "testing"

func TestParseRecordType(t *testing.T) {
	t.Parallel()

	tests := map[string]RecordType{
		"":       RecordA,
		"a":      RecordA,
		" AAAA ": RecordAAAA,
		"cname":  RecordCNAME,
		"mx":     RecordMX,
		"txt":    RecordTXT,
		"ns":     RecordNS,
	}
	for input, expected := range tests {
		actual, err := ParseRecordType(input)
		if err != nil {
			t.Fatalf("ParseRecordType(%q) returned %v", input, err)
		}
		if actual != expected {
			t.Errorf("ParseRecordType(%q) = %q, want %q", input, actual, expected)
		}
	}
}

func TestParseRecordTypeRejectsDeferredTypes(t *testing.T) {
	t.Parallel()

	for _, input := range []string{"SOA", "CAA", "SRV", "unknown"} {
		_, err := ParseRecordType(input)
		if !IsErrorKind(err, ErrUnsupportedType) {
			t.Errorf("ParseRecordType(%q) error = %v, want unsupported-type kind", input, err)
		}
	}
}

func TestSupportedRecordTypesReturnsCopy(t *testing.T) {
	t.Parallel()

	first := SupportedRecordTypes()
	first[0] = RecordSOA
	if SupportedRecordTypes()[0] != RecordA {
		t.Fatal("SupportedRecordTypes returned shared mutable storage")
	}
}
