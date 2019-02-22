package group

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/kafka"
	"github.com/spf13/cobra"
)

var CmdDescGroup = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Short:   "Get Group Details",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var groupMeta []kafkactl.GroupMeta
		switch true {
		case cmd.Flags().Changed("groups"):
			groupMeta = kafka.GroupMetaByTopics(args...)
		default:
			groupMeta = kafka.SearchGroupMeta(args...)
		}
		switch true {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(groupMeta, outFmt))
		default:
			kafka.PrintOut(groupMeta)
		}
	},
}

func init() {
}
