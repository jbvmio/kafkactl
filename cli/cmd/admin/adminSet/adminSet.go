package adminset

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var CmdAdminSet = &cobra.Command{
	Use:   "set",
	Short: "Set Kafka Configurations",
	Run: func(cmd *cobra.Command, args []string) {
		match := true
		switch match {
		case len(args) > 0:
			out.Failf("No such resource: %v", args[0])
		default:
			cmd.Help()
		}
	},
}

func init() {
	//CmdAdminGet.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")

	CmdAdminSet.AddCommand(cmdAdminSetTopic)
	//CmdAdmin.AddCommand(topic.CmdGetTopic)
	//CmdAdmin.AddCommand(group.CmdGetGroup)
	//CmdAdmin.AddCommand(group.CmdGetMember)
	//CmdAdmin.AddCommand(lag.CmdGetLag)
	//CmdAdmin.AddCommand(msg.CmdGetMsg)
}
