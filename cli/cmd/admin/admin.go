package admin

import (
	adminget "github.com/jbvmio/kafkactl/cli/cmd/admin/adminGet"
	adminops "github.com/jbvmio/kafkactl/cli/cmd/admin/adminOps"
	adminset "github.com/jbvmio/kafkactl/cli/cmd/admin/adminSet"
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
		default:
			cmd.Help()
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
	CmdAdmin.AddCommand(cmdAdminRebalance)
}
