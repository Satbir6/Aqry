package dns

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const DefaultTimeout = 3 * time.Second

type LookupRequest struct {
	Domain     string
	RecordType RecordType
	Server     string
	Timeout    time.Duration
	All        bool
}

type Resolver interface {
	Lookup(context.Context, LookupRequest) (LookupResult, error)
}

type SystemResolver struct{}

func NewSystemResolver() *SystemResolver {
	return &SystemResolver{}
}

func (r *SystemResolver) Lookup(ctx context.Context, request LookupRequest) (LookupResult, error) {
	domain, err := NormalizeAndValidateDomain(request.Domain)
	if err != nil {
		return LookupResult{}, err
	}

	recordType, err := ParseRecordType(string(request.RecordType))
	if err != nil {
		return LookupResult{}, err
	}

	timeout := request.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	lookupCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	netResolver, resolverName, err := resolverForServer(request.Server, timeout)
	if err != nil {
		return LookupResult{}, NewLookupError(ErrResolverFailed, domain, recordType, err)
	}

	records, err := lookupRecords(lookupCtx, netResolver, domain, recordType)
	if err != nil {
		return LookupResult{}, classifyLookupError(lookupCtx, domain, recordType, err)
	}
	if len(records) == 0 {
		return LookupResult{}, NewLookupError(ErrNoRecords, domain, recordType, nil)
	}
	if !request.All && len(records) > 1 {
		records = records[:1]
	}

	return LookupResult{
		Domain:     domain,
		RecordType: recordType,
		Resolver:   resolverName,
		Records:    records,
	}, nil
}

func resolverForServer(server string, timeout time.Duration) (*net.Resolver, string, error) {
	server = strings.TrimSpace(server)
	if server == "" || strings.EqualFold(server, "system") || strings.EqualFold(server, "default") {
		return net.DefaultResolver, "system", nil
	}

	address, err := resolverAddress(server)
	if err != nil {
		return nil, "", err
	}

	dialer := &net.Dialer{Timeout: timeout}
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, address)
		},
	}

	return resolver, server, nil
}

func resolverAddress(server string) (string, error) {
	if net.ParseIP(server) != nil {
		return net.JoinHostPort(server, "53"), nil
	}

	host, port, err := net.SplitHostPort(server)
	if err != nil || net.ParseIP(host) == nil {
		return "", fmt.Errorf("resolver must be an IP address, optionally with a port")
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return "", fmt.Errorf("resolver port must be between 1 and 65535")
	}

	return net.JoinHostPort(host, port), nil
}

func ValidateResolverServer(server string) error {
	_, err := resolverAddress(strings.TrimSpace(server))
	return err
}

func lookupRecords(ctx context.Context, resolver *net.Resolver, domain string, recordType RecordType) ([]Record, error) {
	switch recordType {
	case RecordA, RecordAAAA:
		addresses, err := resolver.LookupIPAddr(ctx, domain)
		if err != nil {
			return nil, err
		}
		records := make([]Record, 0, len(addresses))
		seen := make(map[string]struct{}, len(addresses))
		for _, address := range addresses {
			ip := address.IP
			if recordType == RecordA {
				ip = ip.To4()
				if ip == nil {
					continue
				}
			} else if ip.To4() != nil {
				continue
			}
			value := ip.String()
			if _, exists := seen[value]; exists {
				continue
			}
			seen[value] = struct{}{}
			records = append(records, Record{Type: recordType, Name: domain, Value: value})
		}
		return records, nil

	case RecordCNAME:
		canonical, err := resolver.LookupCNAME(ctx, domain)
		if err != nil {
			return nil, err
		}
		canonical = trimDNSRootDot(canonical)
		if strings.EqualFold(canonical, domain) {
			return nil, nil
		}
		return []Record{{Type: recordType, Name: domain, Value: canonical}}, nil

	case RecordMX:
		exchanges, err := resolver.LookupMX(ctx, domain)
		if err != nil {
			return nil, err
		}
		records := make([]Record, 0, len(exchanges))
		for _, exchange := range exchanges {
			records = append(records, Record{
				Type: recordType, Name: domain,
				Value: trimDNSRootDot(exchange.Host), Priority: int(exchange.Pref),
			})
		}
		return records, nil

	case RecordTXT:
		values, err := resolver.LookupTXT(ctx, domain)
		if err != nil {
			return nil, err
		}
		records := make([]Record, 0, len(values))
		for _, value := range values {
			records = append(records, Record{Type: recordType, Name: domain, Value: value})
		}
		return records, nil

	case RecordNS:
		servers, err := resolver.LookupNS(ctx, domain)
		if err != nil {
			return nil, err
		}
		records := make([]Record, 0, len(servers))
		for _, server := range servers {
			records = append(records, Record{
				Type: recordType, Name: domain, Value: trimDNSRootDot(server.Host),
			})
		}
		return records, nil
	default:
		return nil, NewLookupError(ErrUnsupportedType, domain, recordType, nil)
	}
}

func trimDNSRootDot(value string) string {
	if value == "." {
		return value
	}
	return strings.TrimSuffix(value, ".")
}

func classifyLookupError(ctx context.Context, domain string, recordType RecordType, err error) error {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
		return NewLookupError(ErrTimeout, domain, recordType, err)
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		if dnsErr.IsTimeout {
			return NewLookupError(ErrTimeout, domain, recordType, err)
		}
		if dnsErr.IsNotFound {
			return NewLookupError(ErrNoRecords, domain, recordType, err)
		}
	}

	return NewLookupError(ErrResolverFailed, domain, recordType, err)
}
