package cfg

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var cmdView = &cobra.Command{
	Use:     "view",
	Aliases: []string{"show"},
	Short:   "Display kafkactl config",
	Run: func(cmd *cobra.Command, args []string) {
		out.Marshal(GetConfig(), outFlags.Format)
	},
}

func init() {
}
