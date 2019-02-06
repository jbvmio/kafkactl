package adminget

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var CmdAdminGet = &cobra.Command{
	Use:   "get",
	Short: "Get Kafka Configurations",
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
	CmdAdminGet.AddCommand(cmdAdminGetTopic)
	CmdAdminGet.AddCommand(cmdAdminGetPre)
	CmdAdminGet.AddCommand(cmdAdminGetOffsets)
	CmdAdminGet.AddCommand(cmdAdminGetReplicas)
}
