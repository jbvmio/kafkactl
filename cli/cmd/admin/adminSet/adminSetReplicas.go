package adminset

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var cmdAdminSetReplicas = &cobra.Command{
	Use:   "replicas",
	Short: "Set Topic Replicas",
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
	//cmdAdminSetOffsets.AddCommand(cmdAdminGetTopic)
	//cmdAdminSetOffsets.AddCommand(cmdAdminGetPre)
}
