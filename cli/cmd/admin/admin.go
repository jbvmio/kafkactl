package admin

import (
	"github.com/jbvmio/kafkactl/cli/cmd/admin/adminGet"
	"github.com/jbvmio/kafkactl/cli/cmd/admin/adminOps"
	"github.com/jbvmio/kafkactl/cli/cmd/admin/adminSet"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var outFlags out.OutFlags

var CmdAdmin = &cobra.Command{
	Use:   "admin",
	Short: "Kafka Admin Actions",
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case len(args) > 0:
			out.Failf("No such resource: %v", args[0])
		}
	},
}

func init() {
	CmdAdmin.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")

	CmdAdmin.AddCommand(adminget.CmdAdminGet)
	CmdAdmin.AddCommand(adminset.CmdAdminSet)
	CmdAdmin.AddCommand(adminops.CmdAdminCreate)
	CmdAdmin.AddCommand(adminops.CmdAdminDelete)
	CmdAdmin.AddCommand(cmdAdminMove)

	//CmdAdmin.AddCommand(topic.CmdGetTopic)
	//CmdAdmin.AddCommand(group.CmdGetGroup)
	//CmdAdmin.AddCommand(group.CmdGetMember)
	//CmdAdmin.AddCommand(lag.CmdGetLag)
	//CmdAdmin.AddCommand(msg.CmdGetMsg)
}
