package adminops

import (
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var CmdAdminDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete Kafka Resources",
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case len(args) > 0:
			out.Failf("No such resource: %v", args[0])
		default:
			cmd.Help()
		}
	},
}

func init() {
	CmdAdminDelete.AddCommand(cmdAdminDeleteTopic)
	CmdAdminDelete.AddCommand(cmdAdminDeleteGroup)
}
