package cfg

import (
	"github.com/spf13/cobra"
)

var cmdUse = &cobra.Command{
	Use:   "use",
	Short: "Switch between kafkactl contexts",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmdUse.AddCommand(cmdUseContext)
}
