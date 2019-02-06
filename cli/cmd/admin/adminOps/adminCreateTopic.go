package adminops

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var cmdAdminCreateTopic = &cobra.Command{
	Use:   "topic",
	Short: "Create Kafka Topics",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		match := true
		switch match {
		case len(args) > 1:
			out.Warnf("Error: Too many arguments received: %v", args)
			return
		case cmd.Flags().Changed("out"):
			out.Warnf("Error: Cannot use --out when creating topics.")
			return
		default:
			kafka.CreateTopic(args[0], createFlags.PartitionCount, createFlags.ReplicationFactor)
		}
	},
}

func init() {
	cmdAdminCreateTopic.Flags().Int32Var(&createFlags.PartitionCount, "partitions", 0, "Total Number of Partitions for the Topic.")
	cmdAdminCreateTopic.Flags().Int16Var(&createFlags.ReplicationFactor, "replicas", 0, "Total Number of Replicas for the Topic (Replication Factor).")
	cmdAdminCreateTopic.MarkFlagRequired("partitions")
	cmdAdminCreateTopic.MarkFlagRequired("replicas")
}
