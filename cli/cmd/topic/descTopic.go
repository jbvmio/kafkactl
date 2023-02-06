package topic

import (
	"github.com/jbvmio/kafkactl/cli/cmd/lag"
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/kafka"
	"github.com/spf13/cobra"
)

var validOffsets bool

var CmdDescTopic = &cobra.Command{
	Use:     "topic",
	Aliases: []string{"topics"},
	Short:   "Get Topic Details",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lag.ValidOffsets = validOffsets
		var tom []kafkactl.TopicOffsetMap
		var vTom kafka.ValidOffsetTOM
		switch {
		case topicFlags.Lag:
			lag.CmdGetLag.Run(cmd, args)
			return
		case len(topicFlags.Leaders) > 0:
			tom = kafka.FilterTOMByLeader(kafka.SearchTOM(args...), topicFlags.Leaders)
		default:
			tom = kafka.SearchTOM(args...)
		}

		if validOffsets {
			vo := kafka.GetValidOffsets(tom)
			vTom = kafka.ValidOffsetTOM{
				Tom:          tom,
				ValidOffsets: vo,
			}
		}

		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			if validOffsets {
				out.IfErrf(out.Marshal(vTom, outFmt))
				return
			}
			out.IfErrf(out.Marshal(tom, outFmt))
		default:
			if validOffsets {
				kafka.PrintOut(vTom)
				return
			}
			kafka.PrintOut(tom)
		}
	},
}

func init() {
	CmdDescTopic.Flags().BoolVar(&topicFlags.Lag, "lag", false, "Show Any Lag from Specified Topics.")
	CmdDescTopic.Flags().BoolVar(&validOffsets, "valid-offsets", false, "Discover Valid Offsets from Specified Topics.")
	CmdDescTopic.Flags().Int32SliceVar(&topicFlags.Leaders, "leader", []int32{}, "Filter Topic Partitions by Current Leaders")
}
