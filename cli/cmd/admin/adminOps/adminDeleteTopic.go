package adminops

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var cmdAdminDeleteTopic = &cobra.Command{
	Use:   "topic",
	Short: "Delete Kafka Topics",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		match := true
		switch match {
		case cmd.Flags().Changed("out"):
			out.Warnf("Error: Cannot use --out when deleting topics.")
			return
		default:
			kafka.DeleteTopics(args...)
		}
	},
}

func init() {
}
