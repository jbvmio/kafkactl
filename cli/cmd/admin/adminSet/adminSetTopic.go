package adminset

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var topicFlags kafka.TopicConfigFlags

var cmdAdminSetTopic = &cobra.Command{
	Use:       "topic",
	Short:     "Set Topic Configuration",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{"<topic name>"},
	Run: func(cmd *cobra.Command, args []string) {
		var topicConfigs []kafka.TopicConfig
		switch true {
		case topicFlags.SetDefault:
			topicConfigs = kafka.SetDefaultConfig(topicFlags.Config, args...)
		case topicFlags.Config == "" || topicFlags.Value == "":
			out.Warnf("Error: --config and --value flags required unless performing a reset.")
			return
		default:
			topicConfigs = kafka.SetTopicConfig(topicFlags.Config, topicFlags.Value, args...)
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
	cmdAdminSetTopic.Flags().StringVar(&topicFlags.Config, "config", "", "Configuration or Key Names to set.")
	cmdAdminSetTopic.Flags().StringVar(&topicFlags.Value, "value", "", "Configuration Value to set.")
	cmdAdminSetTopic.Flags().BoolVar(&topicFlags.SetDefault, "reset", false, "Set All Topic Configs Back to defaults.")
}
