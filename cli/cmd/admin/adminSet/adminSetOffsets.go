package adminset

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var cmdAdminSetOffsets = &cobra.Command{
	Use:   "offsets",
	Short: "Set Kafka Group Offsets",
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
	//cmdAdminSetOffsets.AddCommand(cmdAdminGetTopic)
	//cmdAdminSetOffsets.AddCommand(cmdAdminGetPre)
}
