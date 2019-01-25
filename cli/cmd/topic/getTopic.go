package topic

import (
	"github.com/jbvmio/kafkactl"
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var topicFlags kafka.TopicFlags

var CmdGetTopic = &cobra.Command{
	Use:     "topic",
	Aliases: []string{"topics"},
	Short:   "Get Topic Info",
	Run: func(cmd *cobra.Command, args []string) {
		var topicSummaries []kafkactl.TopicSummary
		match := true
		switch match {
		case topicFlags.Describe:
			CmdDescTopic.Run(cmd, args)
		default:
			topicSummaries = kafkactl.GetTopicSummaries(kafka.SearchTopicMeta(args...))
		}
		switch match {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.PrintObject(topicSummaries, outFmt)
		default:
			kafka.PrintOut(topicSummaries)
		}
	},
}

func init() {
	CmdGetTopic.Flags().BoolVar(&topicFlags.Describe, "describe", false, "Shortcut/Pass to Describe Command.")
	CmdGetTopic.Flags().StringVar(&topicFlags.Leaders, "leader", "", "Filter Topic Partitions by Current Leaders")
}
