package commands

import (
	"github.com/spf13/cobra"
)

var (
	// migrateCmd represents the init command
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "",
		Long:  ``,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}
