package adminset

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var CmdAdminSet = &cobra.Command{
	Use:   "set",
	Short: "Set Kafka Configurations",
	Run: func(cmd *cobra.Command, args []string) {
		match := true
		switch match {
		case len(args) > 0:
			out.Failf("No such resource: %v", args[0])
		default:
			cmd.Help()
		}
	},
}

func init() {
	CmdAdminSet.AddCommand(cmdAdminSetTopic)
	CmdAdminSet.AddCommand(cmdAdminSetPre)
}
