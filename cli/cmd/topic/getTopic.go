package topic

import (
	"github.com/jbvmio/kafkactl/cli/cmd/group"
	"github.com/jbvmio/kafkactl/cli/cmd/lag"
	"github.com/jbvmio/kafkactl/cli/kafka"
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/kafkactl/kafka"
	"github.com/spf13/cobra"
)

var topicFlags kafka.TopicFlags

var CmdGetTopic = &cobra.Command{
	Use:     "topic",
	Aliases: []string{"topics"},
	Example: examples.GetTopics(),
	Short:   "Get Topic Info",
	Run: func(cmd *cobra.Command, args []string) {
		var topicSummaries []kafkactl.TopicSummary
		switch {
		case topicFlags.Lag:
			lag.CmdGetLag.Run(cmd, args)
			return
		case topicFlags.Group:
			group.CmdDescGroup.Run(cmd, args)
			return
		case topicFlags.Describe || len(topicFlags.Leaders) > 0:
			CmdDescTopic.Run(cmd, args)
			return
		default:
			topicSummaries = kafkactl.GetTopicSummaries(kafka.SearchTopicMeta(args...))
		}
		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(topicSummaries, outFmt))
		default:
			kafka.PrintOut(topicSummaries)
		}
	},
}

func init() {
	CmdGetTopic.Flags().BoolVar(&topicFlags.Describe, "describe", false, "Shortcut/Pass to Describe Command.")
	CmdGetTopic.Flags().BoolVar(&topicFlags.Group, "groups", false, "Show Active Groups Consuming from Specified Topics.")
	CmdGetTopic.Flags().BoolVar(&topicFlags.Lag, "lag", false, "Show Any Lag from Specified Topics.")
	CmdGetTopic.Flags().Int32SliceVar(&topicFlags.Leaders, "leader", []int32{}, "Filter Topic Partitions by Current Leaders")
}
