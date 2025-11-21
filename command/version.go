package command

import (
	"github.com/spf13/cobra"
	"simple-file-server/lib/console"
	"simple-file-server/lib/version"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	RunE: func(cmd *cobra.Command, args []string) error {
		console.GenerateSuccessData(map[string]any{
			"version": version.VERSION,
		})
		return nil
	},
}
