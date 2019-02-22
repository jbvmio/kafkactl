package group

import (
	"github.com/jbvmio/kafkactl/cli/cmd/lag"
	"github.com/jbvmio/kafkactl/cli/kafka"
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/kafka"
	"github.com/spf13/cobra"
)

var groupFlags kafka.GroupFlags

var CmdGetGroup = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Example: examples.GetGroups(),
	Short:   "Get Group Info",
	Run: func(cmd *cobra.Command, args []string) {
		var glm []kafkactl.GroupListMeta
		switch true {
		case groupFlags.Lag:
			lag.CmdGetLag.Run(cmd, args)
			return
		case groupFlags.Describe:
			CmdDescGroup.Run(cmd, args)
			return
		case groupFlags.Topic:
		default:
			glm = kafka.SearchGroupListMeta(args...)
		}
		switch true {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(glm, outFmt))
		default:
			kafka.PrintOut(glm)
		}
	},
}

func init() {
	CmdGetGroup.Flags().BoolVar(&groupFlags.Describe, "describe", false, "Shortcut/Pass to Describe Command.")
	CmdGetGroup.Flags().BoolVar(&groupFlags.Lag, "lag", false, "Shortcut/Pass to Lag Command.")
}
