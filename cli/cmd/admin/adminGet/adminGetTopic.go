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
		switch true {
		default:
			topicConfigs = kafka.SearchTopicConfigs(topicFlags.Configs, args...)
			if topicFlags.GetNonDefaults {
				topicConfigs = kafka.GetNonDefaultConfigs(topicConfigs)
			}
		}
		switch true {
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
	cmdAdminGetTopic.Flags().StringSliceVar(&topicFlags.Configs, "config", []string{}, "Filter by Configuration / Key Names.")
	cmdAdminGetTopic.Flags().BoolVar(&topicFlags.GetNonDefaults, "changed", false, "Only show configs that have changed.")
}
