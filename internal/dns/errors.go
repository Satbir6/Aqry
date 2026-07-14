package dns

import (
	"errors"
	"fmt"
)

type ErrorKind string

const (
	ErrInvalidDomain   ErrorKind = "invalid_domain"
	ErrNoRecords       ErrorKind = "no_records"
	ErrTimeout         ErrorKind = "timeout"
	ErrResolverFailed  ErrorKind = "resolver_failed"
	ErrUnsupportedType ErrorKind = "unsupported_type"
)

type LookupError struct {
	Kind       ErrorKind
	Domain     string
	RecordType RecordType
	Err        error
}

func NewLookupError(kind ErrorKind, domain string, recordType RecordType, err error) *LookupError {
	return &LookupError{Kind: kind, Domain: domain, RecordType: recordType, Err: err}
}

func (e *LookupError) Error() string {
	switch e.Kind {
	case ErrInvalidDomain:
		if e.Err != nil {
			return fmt.Sprintf("invalid domain: %v", e.Err)
		}
		return "invalid domain"
	case ErrNoRecords:
		if e.RecordType != "" && e.Domain != "" {
			return fmt.Sprintf("no %s records found for %s", e.RecordType, e.Domain)
		}
		return "no DNS records found"
	case ErrTimeout:
		return "DNS lookup timed out"
	case ErrResolverFailed:
		if e.Err != nil {
			return fmt.Sprintf("could not reach DNS resolver: %v", e.Err)
		}
		return "could not reach DNS resolver"
	case ErrUnsupportedType:
		if e.RecordType != "" {
			return fmt.Sprintf("record type %s is not supported yet", e.RecordType)
		}
		return "record type is not supported yet"
	default:
		if e.Err != nil {
			return e.Err.Error()
		}
		return "DNS lookup failed"
	}
}

func (e *LookupError) Unwrap() error {
	return e.Err
}

func ErrorKindOf(err error) (ErrorKind, bool) {
	var lookupErr *LookupError
	if errors.As(err, &lookupErr) {
		return lookupErr.Kind, true
	}

	return "", false
}

func IsErrorKind(err error, kind ErrorKind) bool {
	actual, ok := ErrorKindOf(err)
	return ok && actual == kind
}
