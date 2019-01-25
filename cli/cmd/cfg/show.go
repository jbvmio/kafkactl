package cfg

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var cmdShow = &cobra.Command{
	Use:   "show",
	Short: "Display kafkactl configurations",
	Run: func(cmd *cobra.Command, args []string) {
		out.Marshal(GetConfig(), outFlags.Format)
	},
}

func init() {
	cmdShow.AddCommand(cmdShowContext)
}
