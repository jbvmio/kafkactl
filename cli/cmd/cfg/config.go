package cfg

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var outFlags out.OutFlags

var CmdConfig = &cobra.Command{
	Use:   "config",
	Short: "Show and Edit kafkactl config",
}

func init() {
	CmdConfig.PersistentFlags().StringVar(&outFlags.Format, "out", "yaml", "Output Format - yaml|json.")

	CmdConfig.AddCommand(cmdShow)
	CmdConfig.AddCommand(cmdUse)
}
