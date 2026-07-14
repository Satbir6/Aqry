package tui

import "aqry/internal/dns"

type lookupProgressMsg struct {
	id      int
	message string
	percent float64
}

type lookupResponseMsg struct {
	id     int
	result dns.LookupResult
}

type lookupSuccessMsg struct {
	id     int
	result dns.LookupResult
}

type lookupErrorMsg struct {
	id  int
	err error
}

type copySuccessMsg struct{}

type copyErrorMsg struct {
	err error
}
