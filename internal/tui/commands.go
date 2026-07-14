package tui

import (
	"context"
	"fmt"
	"time"

	"aqry/internal/dns"

	tea "github.com/charmbracelet/bubbletea"
)

func lookupCmd(ctx context.Context, id int, resolver dns.Resolver, request dns.LookupRequest) tea.Cmd {
	return func() tea.Msg {
		timeout := request.Timeout
		if timeout <= 0 {
			timeout = dns.DefaultTimeout
		}
		lookupContext, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		result, err := resolver.Lookup(lookupContext, request)
		if err != nil {
			return lookupErrorMsg{id: id, err: err}
		}
		return lookupResponseMsg{id: id, result: result}
	}
}

func progressCmd(id int, delay time.Duration, message string, percent float64) tea.Cmd {
	return tea.Tick(delay, func(time.Time) tea.Msg {
		return lookupProgressMsg{id: id, message: message, percent: percent}
	})
}

func finalizeLookupCmd(id int, result dns.LookupResult) tea.Cmd {
	return func() tea.Msg {
		return lookupSuccessMsg{id: id, result: result}
	}
}

func copyCmd(clipboard Clipboard, value string) tea.Cmd {
	return func() tea.Msg {
		if clipboard == nil {
			return copyErrorMsg{err: fmt.Errorf("clipboard is unavailable")}
		}
		if err := clipboard.WriteAll(value); err != nil {
			return copyErrorMsg{err: err}
		}
		return copySuccessMsg{}
	}
}
