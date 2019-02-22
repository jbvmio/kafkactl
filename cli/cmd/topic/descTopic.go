package topic

import (
	"github.com/jbvmio/kafkactl/cli/cmd/lag"
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/kafka"
	"github.com/spf13/cobra"
)

var CmdDescTopic = &cobra.Command{
	Use:     "topic",
	Aliases: []string{"topics"},
	Short:   "Get Topic Details",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var tom []kafkactl.TopicOffsetMap
		switch true {
		case topicFlags.Lag:
			lag.CmdGetLag.Run(cmd, args)
			return
		case len(topicFlags.Leaders) > 0:
			tom = kafka.FilterTOMByLeader(kafka.SearchTOM(args...), topicFlags.Leaders)
		default:
			tom = kafka.SearchTOM(args...)
		}
		switch true {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(tom, outFmt))
		default:
			kafka.PrintOut(tom)
		}
	},
}

func init() {
	CmdDescTopic.Flags().BoolVar(&topicFlags.Lag, "lag", false, "Show Any Lag from Specified Topics.")
	CmdDescTopic.Flags().Int32SliceVar(&topicFlags.Leaders, "leader", []int32{}, "Filter Topic Partitions by Current Leaders")
}
