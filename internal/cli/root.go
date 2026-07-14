package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"aqry/internal/dns"
	"aqry/internal/tui"
	"aqry/internal/version"

	"github.com/spf13/cobra"
)

type TUIRunner func(tui.Options) error

type Dependencies struct {
	Resolver dns.Resolver
	Out      io.Writer
	Err      io.Writer
	RunTUI   TUIRunner
}

type commandError struct {
	code int
	err  error
}

func (e *commandError) Error() string { return e.err.Error() }
func (e *commandError) Unwrap() error { return e.err }

func NewRootCommand(dependencies Dependencies) *cobra.Command {
	if dependencies.Resolver == nil {
		dependencies.Resolver = dns.NewSystemResolver()
	}
	if dependencies.Out == nil {
		dependencies.Out = os.Stdout
	}
	if dependencies.Err == nil {
		dependencies.Err = os.Stderr
	}
	if dependencies.RunTUI == nil {
		dependencies.RunTUI = tui.Run
	}

	options := flags{}
	command := &cobra.Command{
		Use:           "aqry [domain]",
		Short:         "A fast, interactive DNS resolver",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &commandError{code: 1, err: fmt.Errorf("accepts at most one domain argument")}
			}
			return nil
		},
		RunE: func(command *cobra.Command, args []string) error {
			if options.version {
				_, err := fmt.Fprintf(command.OutOrStdout(), "aqry version %s\n", version.String())
				return err
			}
			if options.timeoutSec <= 0 {
				return &commandError{code: 1, err: fmt.Errorf("timeout must be greater than zero")}
			}

			recordType, err := dns.ParseRecordType(options.recordType)
			if err != nil {
				return err
			}
			timeout := time.Duration(options.timeoutSec) * time.Second

			if options.interactive || len(args) == 0 {
				domain := ""
				if len(args) == 1 {
					domain = args[0]
				}
				if err := dependencies.RunTUI(tui.Options{
					Domain: domain, RecordType: recordType, Server: options.server,
					Timeout: timeout, NoColor: options.noColor, Resolver: dependencies.Resolver,
				}); err != nil {
					return &commandError{code: 5, err: fmt.Errorf("interactive mode failed: %w", err)}
				}
				return nil
			}

			request := dns.LookupRequest{
				Domain: args[0], RecordType: recordType, Server: options.server,
				Timeout: timeout, All: options.all,
			}
			ctx := command.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			result, err := dependencies.Resolver.Lookup(ctx, request)
			if err != nil {
				return err
			}
			return WriteResult(command.OutOrStdout(), result, options.all, options.json)
		},
	}

	command.SetOut(dependencies.Out)
	command.SetErr(dependencies.Err)
	command.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return &commandError{code: 1, err: err}
	})

	flags := command.Flags()
	flags.StringVarP(&options.recordType, "type", "t", "A", "DNS record type (A, AAAA, CNAME, MX, TXT, NS)")
	flags.BoolVar(&options.all, "all", false, "print all records")
	flags.BoolVar(&options.json, "json", false, "print structured JSON")
	flags.StringVarP(&options.server, "server", "s", "system", "DNS resolver IP address")
	flags.IntVar(&options.timeoutSec, "timeout", 3, "lookup timeout in seconds")
	flags.BoolVarP(&options.interactive, "interactive", "i", false, "force interactive mode")
	flags.BoolVar(&options.noColor, "no-color", false, "disable color in interactive mode")
	flags.BoolVarP(&options.version, "version", "v", false, "print version")
	registerCompletions(command)

	return command
}

func Execute() int {
	command := NewRootCommand(Dependencies{})
	err := command.Execute()
	if err == nil {
		return 0
	}
	_, _ = fmt.Fprintf(command.ErrOrStderr(), "aqry: %v\n", err)
	return ExitCode(err)
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var cliErr *commandError
	if errors.As(err, &cliErr) {
		return cliErr.code
	}

	kind, ok := dns.ErrorKindOf(err)
	if !ok {
		return 5
	}
	switch kind {
	case dns.ErrInvalidDomain:
		return 1
	case dns.ErrNoRecords, dns.ErrResolverFailed:
		return 2
	case dns.ErrTimeout:
		return 3
	case dns.ErrUnsupportedType:
		return 4
	default:
		return 5
	}
}
