package send

import (
	"fmt"
	"os"

	"github.com/jbvmio/kafkactl/cli/kafka"
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x"

	"github.com/spf13/cobra"
)

var sendFlags kafka.SendFlags

var CmdSend = &cobra.Command{
	Use:     "send",
	Example: examples.SEND(),
	Aliases: []string{"produce"},
	Short:   "Send/Produce Messages to a Kafka Topic",
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case x.StdinAvailable():
			sendFlags.FromStdin = true
			kafka.ProduceFromFile(sendFlags, os.Stdin, args...)
			return
		case cmd.Flags().Changed("key") || cmd.Flags().Changed("value"):
			kafka.ProduceFromFile(sendFlags, nil, args...)
			return
		default:
			fmt.Println("CONSOLE HERE")
			return
		}
	},
}

func init() {
	CmdSend.Flags().StringVarP(&sendFlags.Delimiter, "delimiter", "D", "", "Delimiter used to separate key,values using Stdin.")
	CmdSend.Flags().StringVarP(&sendFlags.Key, "key", "K", "", "Desired Key for the Message.")
	CmdSend.Flags().StringVarP(&sendFlags.Value, "value", "V", "", "Desired Value for the Message.")
	CmdSend.Flags().BoolVar(&sendFlags.AllPartitions, "allparts", false, "Send Messages to All Available Topic Partitions.")
	CmdSend.Flags().Int32VarP(&sendFlags.Partition, "partition", "p", -1, "Target a Specific Partition, otherwise dictated by partitioning scheme.")
	CmdSend.Flags().StringSliceVar(&sendFlags.Partitions, "partitions", []string{}, "Target Specific Partitions, otherwise dictated by partitioning scheme (comma separated list).")
	CmdSend.Flags().BoolVar(&sendFlags.NoSplit, "no-split", false, "Disable Line Splits.")
	CmdSend.Flags().StringVar(&sendFlags.LineSplit, "split", "\n", "Line Split Delimiter.")
}
