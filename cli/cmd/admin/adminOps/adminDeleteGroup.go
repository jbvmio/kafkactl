package adminops

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var cmdAdminDeleteGroup = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Short:   "Delete Kafka Groups",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case cmd.Flags().Changed("out"):
			out.Warnf("Error: Cannot use --out when deleting groups.")
			return
		default:
			kafka.DeleteGroups(args...)
		}
	},
}

func init() {
}
