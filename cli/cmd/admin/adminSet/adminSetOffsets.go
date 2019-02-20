package adminset

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var offsetFlags kafka.OpsOffsetFlags

var cmdAdminSetOffsets = &cobra.Command{
	Use:     "offsets",
	Aliases: []string{"offset"},
	Short:   "Set Kafka Group Offsets",
	Example: examples.AdminSetOffsets(),
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case cmd.Flags().Changed("out"):
			out.Warnf("Error: Cannot use --out when performing set offset operations.")
			return
		default:
			kafka.SetOffsets(offsetFlags, args...)
		}
	},
}

func init() {
	cmdAdminSetOffsets.Flags().BoolVar(&offsetFlags.AllParts, "allparts", false, "Perform set Offset Operations on all Partitions")
	cmdAdminSetOffsets.Flags().Int32SliceVar(&offsetFlags.Partitions, "partitions", []int32{}, "Partitions to Perform set Offset Operations.")
	cmdAdminSetOffsets.Flags().Int64Var(&offsetFlags.Offset, "offset", 0, "Target offset.")
	cmdAdminSetOffsets.Flags().StringVar(&offsetFlags.Group, "group", "", "Target Group to perform set Offset Operation.")
	cmdAdminSetOffsets.Flags().BoolVar(&offsetFlags.OffsetNewest, "newest", false, "Target Newest Offset available on Topic.")
	cmdAdminSetOffsets.Flags().BoolVar(&offsetFlags.OffsetOldest, "oldest", false, "Target Oldest Offset available on Topic.")
	cmdAdminSetOffsets.Flags().Int64Var(&offsetFlags.RelativeOffset, "relative", 0, "Relative Offset to Target from Newest.")
}
