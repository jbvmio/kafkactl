package zk

import (
	"github.com/jbvmio/kafkactl/cli/cmd/cfg"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/jbvmio/kafkactl/cli/zookeeper"
	"github.com/spf13/cobra"
)

var cxFlags cfg.CXFlags
var outFlags out.OutFlags
var zkFlags zookeeper.ZKFlags

var CmdZK = &cobra.Command{
	Use:   "zk",
	Short: "Zookeeper Actions",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		switch {
		case cmd.Flags().Changed("zookeeper"):
			zookeeper.LaunchZKClient(cfg.AdhocContext(cxFlags), zkFlags)
		default:
			zookeeper.LaunchZKClient(cfg.GetContext(zkFlags.Context), zkFlags)
		}
	},
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
	CmdZK.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")
	CmdZK.PersistentFlags().StringVarP(&zkFlags.Context, "context", "C", "", "Specify a context.")
	CmdZK.PersistentFlags().BoolVarP(&zkFlags.Verbose, "verbose", "v", false, "Display additional info or errors.")
	CmdZK.PersistentFlags().StringVarP(&cxFlags.Broker, "broker", "B", "", "Specify a single broker target host:port - Overrides config.")
	CmdZK.PersistentFlags().StringVarP(&cxFlags.Zookeeper, "zookeeper", "Z", "", "Specify a single zookeeper target host:port - Overrides config.")

	CmdZK.AddCommand(cmdZKls)
	CmdZK.AddCommand(cmdZKcreate)
	CmdZK.AddCommand(cmdZKdelete)
}
