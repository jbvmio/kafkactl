package zk

import (
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/jbvmio/kafkactl/cli/zookeeper"

	"github.com/spf13/cobra"
)

var cmdZKls = &cobra.Command{
	Use:     "ls",
	Short:   "Print Zookeeper Paths and Values",
	Example: examples.ZKLS(),
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var zkpv []zookeeper.ZKPathValue
		var zkp []zookeeper.ZKPath
		switch true {
		case zkFlags.Recurse:
			zkp = zookeeper.ZKRecurseLS(zkFlags.Depth, args...)
			if zkFlags.Values {
				zkp = zookeeper.ZKFilterAllVals(zkp)
			}
		default:
			zkpv = zookeeper.ZKls(args...)
		}
		switch true {
		case outFlags.Format != "":
			if zkFlags.Recurse {
				out.IfErrf(out.Marshal(zkp, outFlags.Format))
				return
			}
			out.IfErrf(out.Marshal(zkpv, outFlags.Format))
		default:
			if zkFlags.Recurse {
				printZK(zkp)
				return
			}
			printZK(zkpv)
		}
	},
}

func init() {
	cmdZKls.Flags().BoolVarP(&zkFlags.Recurse, "recurse", "R", false, "List Recursively")
	cmdZKls.Flags().BoolVar(&zkFlags.Values, "values", false, "Return non-empty Values Only (Used with recurse)")
	cmdZKls.Flags().Uint8VarP(&zkFlags.Depth, "depth", "D", 3, "Specify Recursive Depth")
}
