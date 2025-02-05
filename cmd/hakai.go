package cmd

import "github.com/spf13/cobra"

func help(cmd *cobra.Command, _ []string) {
	cmd.Help()
}

func newHakaiCmd() *cobra.Command {
	hakaiCmd := &cobra.Command{
		Use:   "hakai",
		Short: "Hakai",
		Run:   help,
		RunE:  func(cmd *cobra.Command, args []string) error {
            return nil
        },
	}

	return hakaiCmd
}
