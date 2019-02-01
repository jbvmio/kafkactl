package adminget

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var cmdAdminGetReplicas = &cobra.Command{
	Use:   "replicas",
	Short: "Get Topic Replicas",
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
	//cmdAdminGetReplicas.AddCommand(cmdAdminGetTopic)
	//cmdAdminGetReplicas.AddCommand(cmdAdminGetPre)
}
