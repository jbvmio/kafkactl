package bur

import (
	jbbur "github.com/jbvmio/burrow"
	"github.com/jbvmio/kafkactl/cli/burrow"
	"github.com/jbvmio/kafkactl/cli/cmd/cfg"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var cxFlags cfg.CXFlags
var outFlags out.OutFlags
var burFlags burrow.BurrowFlags

var CmdBur = &cobra.Command{
	Use:     "burrow",
	Example: "  kafkactl burrow <groupName>",
	Short:   "Show Burrow Lag Evaluations <wip>",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		switch {
		case cmd.Flags().Changed("burrow"):
			burrow.LaunchBurrowClient(cfg.AdhocContext(cxFlags), burFlags)
		default:
			burrow.LaunchBurrowClient(cfg.GetContext(burFlags.Context), burFlags)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var burrowParts []jbbur.Partition
		switch {
		default:
			burrowParts = burrow.SearchBurrowConsumers(burFlags, args...)
		}
		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(burrowParts, outFmt))
		default:
			printBur(burrowParts)
		}
	},
}

func init() {
	CmdBur.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")
	CmdBur.PersistentFlags().StringVarP(&burFlags.Context, "context", "C", "", "Specify a context.")
	CmdBur.PersistentFlags().StringVar(&cxFlags.Burrow, "burrow", "", "Specify a single burrow endpoint http://host:port - Overrides config.")
	CmdBur.PersistentFlags().BoolVarP(&burFlags.Exact, "exact", "x", false, "Find exact matches.")
	CmdBur.Flags().BoolVar(&burFlags.ErrOnly, "errs", false, "Filter for NON OK status.")
	CmdBur.Flags().StringVarP(&burFlags.Topic, "topic", "t", "", "Filter by a Topic.")

	CmdBur.AddCommand(cmdBurMon)
}
