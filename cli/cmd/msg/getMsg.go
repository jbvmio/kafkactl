package msg

import (
	"github.com/spf13/cobra"
)

var CmdGetMsg = &cobra.Command{
	Use:     "msg",
	Aliases: []string{"log", "logs", "msgs"},
	Short:   "Shortcut to logs/msgs command",
	Args:    cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		logsFlags.TailTouched = cmd.Flags().Changed("tail")
	},
	Run: func(cmd *cobra.Command, args []string) {
		match := true
		switch match {
		default:
			CmdLogs.Run(cmd, args)
			return
		}
	},
}

func init() {
	CmdGetMsg.Flags().BoolVar(&outFlags.Header, "no-header", false, "Suppress Header Information.")
	CmdGetMsg.Flags().BoolVar(&logsFlags.Follow, "follow", false, "Output messages as they Arrive for the Given Topic and Partitions.")
	CmdGetMsg.Flags().Int64Var(&logsFlags.Tail, "tail", 1, "Messages Starting from this Relative Offset to Retrieve.")
	CmdGetMsg.Flags().Int64Var(&logsFlags.Offset, "offset", -1, "Target a Specific Offset.")
	CmdGetMsg.Flags().Int32VarP(&logsFlags.Partition, "partition", "p", -1, "Target a Specific Partition, otherwise all.")
	CmdGetMsg.Flags().StringSliceVar(&logsFlags.Partitions, "partitions", []string{}, "Target Specific Partitions, otherwise all (comma separated list).")
	CmdGetMsg.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")
}
