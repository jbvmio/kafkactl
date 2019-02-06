package zk

import (
	"github.com/jbvmio/kafkactl/cli/zookeeper"
	"github.com/spf13/cobra"
)

var cmdZKcreate = &cobra.Command{
	Use:   "create",
	Short: "Create Zookeeper Paths and Values",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case len([]byte(zkFlags.Value)) > 0:
			zookeeper.ZKCreate(args[0], false, zkFlags.Force, []byte(zkFlags.Value)...)
		default:
			zookeeper.ZKCreate(args[0], false, zkFlags.Force)
		}
	},
}

func init() {
	cmdZKcreate.Flags().BoolVarP(&zkFlags.Force, "force", "F", false, "Force Operation / Create Parent Paths if Necessary")
	cmdZKcreate.Flags().StringVar(&zkFlags.Value, "value", "", "Value to Create or Set")
}
