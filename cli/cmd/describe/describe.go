package describe

import (
	"github.com/jbvmio/kafkactl/cli/cmd/broker"
	"github.com/jbvmio/kafkactl/cli/cmd/group"
	"github.com/jbvmio/kafkactl/cli/cmd/topic"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var outFlags out.OutFlags

var CmdDescribe = &cobra.Command{
	Use:     "describe",
	Aliases: []string{"desc"},
	Short:   "Get Kafka Details",
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
	CmdDescribe.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")

	CmdDescribe.AddCommand(broker.CmdGetBroker)
	CmdDescribe.AddCommand(topic.CmdDescTopic)
	CmdDescribe.AddCommand(group.CmdDescGroup)
}
