package admin

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var rebalanceFlags kafka.OpsReplicaFlags

var cmdAdminRebalance = &cobra.Command{
	Use:   "rebalance",
	Short: "Rebalance Topics Across Brokers",
	//Example: examples.AdminMoveFunc(),
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var rapList kafka.RAPartList
		switch {
		default:
			rapList = kafka.RebalanceTopics(rebalanceFlags, args...)
			if !rebalanceFlags.DryRun {
				ok := kafka.ZkCreateRAP(rapList)
				if !ok {
					out.Warnf("Error Creating Reassign Partitions.")
					return
				}
			} else {
				out.Infof("Performing Dry Run ...")
			}
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
	cmdAdminRebalance.Flags().Int32SliceVar(&rebalanceFlags.Brokers, "brokers", []int32{}, "Desired Broker Leaders.")
	cmdAdminRebalance.Flags().BoolVar(&rebalanceFlags.PreserveLeader, "preserve-leaders", false, "Keep the current partition leaders in place, balance only peer replicas.")
	cmdAdminRebalance.Flags().BoolVar(&rebalanceFlags.DryRun, "dry-run", false, "Perform a Dry Run.")
}
