package group

import (
	"github.com/jbvmio/kafkactl"
	"github.com/jbvmio/kafkactl/cli/cmd/lag"
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var memberFlags kafka.GroupFlags

var CmdGetMember = &cobra.Command{
	Use:     "member",
	Example: `  kafkactl get member <memberID> --lag`,
	Aliases: []string{"clientid"},
	Short:   "Get Groups by Member IDs",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var groupMeta []kafkactl.GroupMeta
		switch true {
		case memberFlags.Lag:
			lag.CmdGetLag.Run(cmd, args)
			return
		default:
			groupMeta = kafka.GroupMetaByMember(args...)
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
	CmdGetMember.Flags().BoolVar(&memberFlags.Lag, "lag", false, "Shortcut/Pass to Lag Command.")
}
