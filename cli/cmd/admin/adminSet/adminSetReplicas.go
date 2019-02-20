package adminset

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var replicaFlags kafka.OpsReplicaFlags

var cmdAdminSetReplicas = &cobra.Command{
	Use:     "replicas",
	Aliases: []string{"replica"},
	Short:   "Set Topic Replicas",
	Example: examples.AdminSetReplicas(),
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var rapList kafka.RAPartList
		switch {
		default:
			rapList = kafka.SetTopicReplicas(replicaFlags, args...)
			if !replicaFlags.DryRun {
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
	cmdAdminSetReplicas.Flags().BoolVar(&replicaFlags.AllParts, "allparts", false, "Target all Partitions.")
	cmdAdminSetReplicas.Flags().Int32SliceVar(&replicaFlags.Brokers, "brokers", []int32{}, "Desired Brokers.")
	cmdAdminSetReplicas.Flags().Int32SliceVar(&replicaFlags.Partitions, "partitions", []int32{}, "Target Partitions.")
	cmdAdminSetReplicas.Flags().IntVar(&replicaFlags.ReplicationFactor, "replicas", 0, "Desired Replication Factor.")
	cmdAdminSetReplicas.Flags().BoolVar(&replicaFlags.DryRun, "dry-run", false, "Perform a Dry Run.")
}
