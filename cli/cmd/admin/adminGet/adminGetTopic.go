package adminget

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var topicFlags kafka.TopicConfigFlags

var cmdAdminGetTopic = &cobra.Command{
	Use:   "topic",
	Short: "Get Topic Configurations",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var topicConfigs []kafka.TopicConfig
		match := true
		switch match {
		default:
			topicConfigs = kafka.SearchTopicConfigs(topicFlags.Configs, args...)
		}
		switch match {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(topicConfigs, outFmt))
		default:
			kafka.PrintOut(topicConfigs)
		}
	},
}

func init() {
	//cmdAdminGetTopic.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")
	cmdAdminGetTopic.Flags().StringSliceVar(&topicFlags.Configs, "configs", []string{}, "Configuration or Key Names.")

	//CmdAdmin.AddCommand(broker.CmdGetBroker)
	//CmdAdmin.AddCommand(topic.CmdGetTopic)
	//CmdAdmin.AddCommand(group.CmdGetGroup)
	//CmdAdmin.AddCommand(group.CmdGetMember)
	//CmdAdmin.AddCommand(lag.CmdGetLag)
	//CmdAdmin.AddCommand(msg.CmdGetMsg)
}
