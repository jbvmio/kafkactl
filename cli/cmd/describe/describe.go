package describe

import (
	"github.com/jbvmio/kafkactl/cli/cmd/group"
	"github.com/jbvmio/kafkactl/cli/cmd/topic"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var outFlags out.OutFlags

var CmdDescribe = &cobra.Command{
	Use:   "describe",
	Short: "Get Kafka Details",
}

func init() {
	CmdDescribe.PersistentFlags().StringVar(&outFlags.Format, "out", "", "Change Output Format - yaml|json.")

	CmdDescribe.AddCommand(topic.CmdDescTopic)
	CmdDescribe.AddCommand(group.CmdDescGroup)
}
