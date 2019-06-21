package msg

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var CmdQuery = &cobra.Command{
	Use: "query",
	//Example: `  kafkactl get msg -p 1 <topicName>`,
	Aliases: []string{"q"},
	Short:   "Query a Topic for results",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryMap kafka.OffsetRangeMap
		switch {
		default:
			queryMap = kafka.GetMsgOffsets(logsFlags, args...)
		}

		switch {
		case cmd.Flags().Changed("out"):
			out.IfErrf(out.Marshal(queryMap, outFlags.Format))
		default:
			kafka.PrintOut(queryMap)
		}

	},
}

func init() {
	CmdQuery.Flags().StringVar(&logsFlags.FromTime, "from", "", `Specify a Starting datetime for the query (Ex: "11/12/2018 03:59:38.508").`)
	CmdQuery.Flags().StringVar(&logsFlags.ToTime, "to", "", `Specify an Ending datetime for the query. (Ex: "11/12/2018 03:59:38.508")`)
	CmdQuery.Flags().StringVar(&logsFlags.LastDuration, "last", "", `Specify the last "x" duration for the query (Ex: "1m").`)
	CmdQuery.Flags().BoolVar(&outFlags.Header, "no-header", false, "Suppress Header Information.")
	CmdQuery.Flags().Int32VarP(&logsFlags.Partition, "partition", "p", -1, "Target a Specific Partition, otherwise all.")
	CmdQuery.Flags().StringSliceVar(&logsFlags.Partitions, "partitions", []string{}, "Target Specific Partitions, otherwise all (comma separated list).")
	CmdQuery.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")
}
