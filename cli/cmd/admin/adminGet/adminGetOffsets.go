package adminget

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var cmdAdminGetOffsets = &cobra.Command{
	Use:   "offsets",
	Short: "Get Kafka Group Offsets",
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
	//cmdAdminGetOffsets.AddCommand(cmdAdminGetTopic)
	//cmdAdminGetOffsets.AddCommand(cmdAdminGetPre)
}
