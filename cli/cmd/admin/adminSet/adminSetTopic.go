package adminset

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var topicFlags kafka.TopicConfigFlags

var cmdAdminSetTopic = &cobra.Command{
	Use:   "topic",
	Short: "Set Topic Configuration",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var topicConfigs []kafka.TopicConfig
		match := true
		switch match {
		default:
			topicConfigs = kafka.SetTopicConfig(topicFlags.Config, topicFlags.Value, args...)
		}
		switch match {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(topicConfigs, outFmt))
		default:
			kafka.PrintAdm(topicConfigs)
		}
	},
}

func init() {
	cmdAdminSetTopic.Flags().StringVar(&topicFlags.Config, "config", "", "Configuration or Key Names to set.")
	cmdAdminSetTopic.Flags().StringVar(&topicFlags.Value, "value", "", "Configuration Value to set.")
	cmdAdminSetTopic.MarkFlagRequired("config")
	cmdAdminSetTopic.MarkFlagRequired("value")
}
