package topic

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var CmdDescTopic = &cobra.Command{
	Use:     "topic",
	Aliases: []string{"topics"},
	Short:   "Get Topic Details",
	Run: func(cmd *cobra.Command, args []string) {
		match := true
		switch match {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.PrintObject(kafka.SearchTOM(args...), outFmt)
		default:
			kafka.PrintOut(kafka.SearchTOM(args...))
		}
	},
}

func init() {
}
