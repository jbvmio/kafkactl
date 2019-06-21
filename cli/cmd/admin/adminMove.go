package admin

import (
	"os"

	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x"
	"github.com/jbvmio/kafkactl/cli/x/out"

	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"

	"github.com/spf13/cobra"
)

var replicaFlags kafka.OpsReplicaFlags

var cmdAdminMove = &cobra.Command{
	Use:     "move",
	Short:   "Move Partitions using Stdin",
	Example: examples.AdminMoveFunc(),
	Args:    cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var rapList kafka.RAPartList
		switch {
		case x.StdinAvailable():
			rapList = kafka.MovePartitionsStdin(kafka.ParseTopicStdin(os.Stdin), replicaFlags.Brokers)
			if !replicaFlags.DryRun {
				ok := kafka.ZkCreateRAP(rapList)
				if !ok {
					out.Warnf("Error Creating Reassign Partitions.")
					return
				}
			} else {
				out.Infof("Performing Dry Run ...")
			}
		default:
			out.Warnf("Command requires piped stdin data.")
			return
		}
		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(rapList, outFmt))
		default:
			kafka.PrintAdm(rapList)
		}
	},
}

func init() {
	cmdAdminMove.Flags().Int32SliceVar(&replicaFlags.Brokers, "brokers", []int32{}, "Desired Brokers.")
	cmdAdminMove.Flags().BoolVar(&replicaFlags.DryRun, "dry-run", false, "Perform a Dry Run.")
	cmdAdminMove.MarkFlagRequired("brokers")
	cmdAdminMove.Flags().BoolVar(&kafka.FORCE, "Force", false, "bypasses any configured checks.")
	cmdAdminMove.Flags().MarkHidden("Force")
}
