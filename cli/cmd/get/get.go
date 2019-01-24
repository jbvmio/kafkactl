package get

import (
	"github.com/jbvmio/kafkactl/cli/cmd/broker"
	"github.com/jbvmio/kafkactl/cli/cmd/group"
	"github.com/jbvmio/kafkactl/cli/cmd/topic"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var outFlags out.OutFlags

var CmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get Kafka Information",
}

func init() {
	CmdGet.PersistentFlags().StringVar(&outFlags.Format, "out", "", "Change Output Format - yaml|json.")

	CmdGet.AddCommand(topic.CmdGetTopic)
	CmdGet.AddCommand(group.CmdGetGroup)
	CmdGet.AddCommand(broker.CmdGetBroker)
}
