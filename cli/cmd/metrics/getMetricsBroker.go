package metrics

import (
	kafkactl "github.com/jbvmio/kafka"
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var metricFlags kafka.MetricFlags

var CmdGetMetrics = &cobra.Command{
	Use:     "metrics",
	Aliases: []string{"metric"},
	Short:   "Get Metrics",
	Run: func(cmd *cobra.Command, args []string) {
		var MC kafkactl.MetricCollection
		switch {
		case metricFlags.Broker:
			MC = kafka.GetKakfaMetrics(metricFlags)
		default:
			cmd.Help()
			return
		}
		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(MC, outFmt))
		default:
			kafka.PrintMetricCollection(MC)
		}
	},
}

func init() {
	CmdGetMetrics.Flags().BoolVar(&metricFlags.Broker, "brokers", false, "Get Broker Metrics.")
	CmdGetMetrics.Flags().IntVar(&metricFlags.Intervals, "intervals", 10, "Number of Intervals to Conduct.")
	CmdGetMetrics.Flags().IntVar(&metricFlags.Seconds, "seconds", 1, "Seconds Between each Interval.")
}
