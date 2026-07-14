package cli

import (
	"fmt"

	"aqry/internal/dns"

	"github.com/spf13/cobra"
)

func registerCompletions(command *cobra.Command) {
	_ = command.RegisterFlagCompletionFunc("type", func(
		_ *cobra.Command, _ []string, _ string,
	) ([]string, cobra.ShellCompDirective) {
		values := make([]string, 0, len(dns.SupportedRecordTypes()))
		for _, recordType := range dns.SupportedRecordTypes() {
			values = append(values, fmt.Sprintf("%s\t%s", recordType, recordType.Description()))
		}
		return values, cobra.ShellCompDirectiveNoFileComp
	})
}
