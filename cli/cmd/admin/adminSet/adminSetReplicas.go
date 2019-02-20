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
			ok := kafka.ZkCreateRAP(rapList)
			if !ok {
				out.Warnf("Error Creating Reassign Partitions.")
				return
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

	//cmdAdminSetOffsets.AddCommand(cmdAdminGetTopic)
	//cmdAdminSetOffsets.AddCommand(cmdAdminGetPre)
}
