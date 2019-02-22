package adminops

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var cmdAdminCreatePartitions = &cobra.Command{
	Use:   "partitions",
	Short: "Create Kafka Topic Partitions",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case cmd.Flags().Changed("out"):
			out.Warnf("Error: Cannot use --out when increasing partitions.")
			return
		default:
			kafka.ConfigurePartitionCount(createFlags, args...)
		}
	},
}

func init() {
	cmdAdminCreatePartitions.Flags().Int32Var(&createFlags.PartitionCount, "partitions", 0, "Total Number of Partitions Desired for the Topic.")
	cmdAdminCreatePartitions.Flags().BoolVar(&createFlags.DryRun, "dry-run", false, "Perform a Dry Run.")
	cmdAdminCreatePartitions.MarkFlagRequired("partitions")
}
