package topic

import (
	"github.com/jbvmio/kafkactl"
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var CmdGetTopicConfig = &cobra.Command{
	Use:   "config",
	Short: "Get Topic Configuration Details",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var tom []kafkactl.TopicOffsetMap
		match := true
		switch match {
		default:
			tom = kafka.SearchTOM(args...)
		}
		switch match {
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
}
