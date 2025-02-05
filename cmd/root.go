package cmd

import (
	"fmt"

	"github.com/lucasmendesl/beerus/version"
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command.
//
// It creates a new root command with the default run, help, and version commands.
// The root command is configured with a custom help and version template,
// and the default command is set to 'help'.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "beerus",
		Short:             shortDescription,
		Long:              longDescription,
		Run:               help,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Version: version.Version,
	}

	rootCmd.SetHelpTemplate(helpTemplate)
	rootCmd.SetVersionTemplate(fmt.Sprintf("Beerus version {{.Version}}, release %s, build %s\n", version.ReleaseDate, version.GitCommit))

	rootCmd.AddCommand(newHakaiCmd())
	rootCmd.PersistentFlags().String("config-file", "", "config file (default is $HOME/.beerus.yaml)")

	return rootCmd
}
