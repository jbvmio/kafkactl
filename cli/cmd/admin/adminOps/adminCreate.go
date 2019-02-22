package adminops

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var createFlags kafka.OpsCreateFlags

var CmdAdminCreate = &cobra.Command{
	Use:   "create",
	Short: "Create Kafka Resources",
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
	CmdAdminCreate.AddCommand(cmdAdminCreateTopic)
	CmdAdminCreate.AddCommand(cmdAdminCreatePartitions)
}
