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
			if topicFlags.GetNonDefaults {
				topicConfigs = kafka.GetNonDefaultConfigs(topicConfigs)
			}
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
	cmdAdminGetTopic.Flags().StringSliceVar(&topicFlags.Configs, "filter", []string{}, "Filter by Configuration / Key Names.")
	cmdAdminGetTopic.Flags().BoolVar(&topicFlags.GetNonDefaults, "non-defaults", false, "Only show configs not using broker defaults.")
}
