package msg

import (
	kafkactl "github.com/jbvmio/kafka"

	"github.com/jbvmio/kafkactl/cli/kafka"
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var logsFlags kafka.MSGFlags
var outFlags out.OutFlags

var CmdLogs = &cobra.Command{
	Use:     "logs",
	Example: examples.LOGS(),
	Aliases: []string{"consume"},
	Short:   "Get Messages from a Kafka Topic",
	Args:    cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		logsFlags.TailTouched = cmd.Flags().Changed("tail")
	},
	Run: func(cmd *cobra.Command, args []string) {
		var msgs []*kafkactl.Message
		switch {
		case logsFlags.Follow:
			kafka.FollowTopic(logsFlags, outFlags, args...)
			return
		default:
			msgs = kafka.GetMessages(logsFlags, args...)
		}
		switch {
		case cmd.Flags().Changed("out"):
			out.IfErrf(out.Marshal(msgs, outFlags.Format))
		default:
			kafka.PrintMSGs(msgs, outFlags)
		}
	},
}

func init() {
	CmdLogs.Flags().BoolVar(&outFlags.Header, "no-header", false, "Suppress Header Information.")
	CmdLogs.Flags().BoolVar(&logsFlags.Follow, "follow", false, "Output messages as they Arrive for the Given Topic and Partitions.")
	CmdLogs.Flags().Int64Var(&logsFlags.Tail, "tail", 1, "Relative Value back from the Newest Offset to Start Message Retrieval.")
	CmdLogs.Flags().Int64Var(&logsFlags.Offset, "offset", -1, "Target a Specific Offset.")
	CmdLogs.Flags().Int32VarP(&logsFlags.Partition, "partition", "p", -1, "Target a Specific Partition, otherwise all.")
	CmdLogs.Flags().StringSliceVar(&logsFlags.Partitions, "partitions", []string{}, "Target Specific Partitions, otherwise all (comma separated list).")
	CmdLogs.Flags().StringSliceVar(&logsFlags.JSONFilters, "json-filter", []string{}, "Filter Message Stream by JSON Filter, used with --follow.")
	CmdLogs.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")
	CmdLogs.AddCommand(CmdQuery)
}
