package zk

import (
	"github.com/jbvmio/kafkactl/cli/zookeeper"
	"github.com/spf13/cobra"
)

var cmdZKdelete = &cobra.Command{
	Use:     "delete",
	Short:   "Delete Zookeeper Paths and Values",
	Example: `  kafkactl zk delete /path/to/delete`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case len(args) == 1:
			zookeeper.ZKDelete(args[0], zkFlags.Recurse)
		default:
			cmd.Help()
		}
	},
}

func init() {
	cmdZKdelete.Flags().BoolVar(&zkFlags.Recurse, "RMR", false, "Delete Recursively - All Child Subdirectories")
}
