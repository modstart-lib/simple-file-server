package command

import (
	"github.com/spf13/cobra"
	"simple-file-server/server"
)

func init() {}

var RootCmd = &cobra.Command{
	Use:   "simple-file-server",
	Short: "simple file server , support upload ( multipart ), url download",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.Start()
		return nil
	},
}
