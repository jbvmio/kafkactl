package broker

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var getMetrics bool

var CmdGetBroker = &cobra.Command{
	Use:     "broker",
	Aliases: []string{"brokers"},
	Short:   "Get Broker Details",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(kafka.GetBrokerInfo(args...), outFmt))
		case getMetrics:

			kafka.TestMetrics()
			/*
				metricRegistry := metrics.NewRegistry()
				histogram := getOrRegisterHistogram("name", metricRegistry)

				if histogram == nil {
					t.Error("Unexpected nil histogram")
				}

				// Fetch the metric
				foundHistogram := metricRegistry.Get("name")

				if foundHistogram != histogram {
					t.Error("Unexpected different histogram", foundHistogram, histogram)
				}

				// Try to register the metric again
				sameHistogram := getOrRegisterHistogram("name", metricRegistry)

				if sameHistogram != histogram {
					t.Error("Unexpected different histogram", sameHistogram, histogram)
				}
			*/
			return
		default:
			kafka.PrintOut(kafka.GetBrokerInfo(args...))
		}
	},
}

func init() {
	CmdGetBroker.Flags().BoolVar(&getMetrics, "metrics", false, "Get Broker Metrics.")
}
