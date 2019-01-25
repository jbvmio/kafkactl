package broker

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var CmdGetBroker = &cobra.Command{
	Use:     "broker",
	Aliases: []string{"brokers"},
	Short:   "Get Broker Details",
	Run: func(cmd *cobra.Command, args []string) {
		match := true
		switch match {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.Marshal(kafka.GetBrokerInfo(args...), outFmt)
		default:
			kafka.PrintOut(kafka.GetBrokerInfo(args...))
		}
	},
}

func init() {
}
