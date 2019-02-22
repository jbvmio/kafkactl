package bur

import (
	"github.com/jbvmio/kafkactl/cli/burrow"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var cmdBurMon = &cobra.Command{
	Use:     "monitor",
	Example: "  kafkactl burrow monitor --group <groupName> --topic <topicName>",
	Short:   "Monitor Lag for a Group and Topic using Terminal Graphs",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case cmd.Flags().Changed("out"):
			out.Failf("Cannot use --out with burrow monitor.")
		default:
			burrow.LaunchBurrowMonitor(burFlags)
		}
	},
}

func init() {
	cmdBurMon.Flags().StringVarP(&burFlags.Group, "group", "g", "", "Group to Monitor.")
	cmdBurMon.Flags().StringVarP(&burFlags.Topic, "topic", "t", "", "Topic to Monitor.")
	cmdBurMon.MarkFlagRequired("group")
	cmdBurMon.MarkFlagRequired("topic")
}
