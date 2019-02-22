package adminget

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"

	"github.com/spf13/cobra"
)

var offsetFlags kafka.OpsOffsetFlags

var cmdAdminGetOffsets = &cobra.Command{
	Use:     "offsets",
	Short:   "Get Kafka Group Offsets",
	Example: examples.AdminGetOffsets(),
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var offsetDetails kafka.OffsetDetails
		switch {
		case offsetFlags.ShowGroups:
			offsetDetails = kafka.MatchGroupOffsets(kafka.GetTopicOffsets(args...))
		default:
			offsetDetails = kafka.GetTopicOffsets(args...)
		}
		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(offsetDetails, outFmt))
		default:
			kafka.PrintAdm(offsetDetails)
		}
	},
}

func init() {
	cmdAdminGetOffsets.Flags().BoolVar(&offsetFlags.ShowGroups, "groups", false, "Show Group Offsets in Addtions to Topic Offsets.")
}
