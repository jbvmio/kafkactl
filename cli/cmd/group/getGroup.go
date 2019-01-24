package group

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var CmdGetGroup = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Short:   "Get Group Details",
	Run: func(cmd *cobra.Command, args []string) {
		match := true
		switch match {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.PrintObject(kafka.SearchGroupListMeta(args...), outFmt)
		default:
			kafka.PrintOut(kafka.SearchGroupListMeta(args...))
		}
	},
}

func init() {
}
