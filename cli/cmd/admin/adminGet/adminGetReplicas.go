package adminget

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var replicaFlags kafka.OpsReplicaFlags
var noSim bool

var cmdAdminGetReplicas = &cobra.Command{
	Use:     "replicas",
	Aliases: []string{"replica"},
	Short:   "Get Topic Replicas or Simulate Replica Set Operations.",
	Example: examples.AdminGetReplicas(),
	Run: func(cmd *cobra.Command, args []string) {
		var rapList kafka.RAPartList
		var replicaDetails kafka.ReplicaDetails
		switch {
		case len(replicaFlags.Brokers) > 0 || len(replicaFlags.Partitions) > 0 || replicaFlags.AllParts || replicaFlags.ReplicationFactor != 0:
			rapList = kafka.SetTopicReplicas(replicaFlags, args...)
		default:
			replicaDetails = kafka.GetTopicReplicas(args...)
			noSim = true
		}
		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			if noSim {
				out.IfErrf(out.Marshal(replicaDetails, outFmt))
				return
			}
			out.IfErrf(out.Marshal(rapList, outFmt))
		default:
			if noSim {
				kafka.PrintAdm(replicaDetails)
				return
			}
			kafka.PrintAdm(rapList)
		}
	},
}

func init() {
	cmdAdminGetReplicas.Flags().BoolVar(&replicaFlags.AllParts, "allparts", false, "Target all Partitions.")
	cmdAdminGetReplicas.Flags().Int32SliceVar(&replicaFlags.Brokers, "brokers", []int32{}, "Desired Brokers.")
	cmdAdminGetReplicas.Flags().Int32SliceVar(&replicaFlags.Partitions, "partitions", []int32{}, "Target Partitions.")
	cmdAdminGetReplicas.Flags().IntVar(&replicaFlags.ReplicationFactor, "replicas", 0, "Desired Replication Factor.")

}
