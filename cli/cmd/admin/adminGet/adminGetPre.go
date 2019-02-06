package adminget

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var preFlags kafka.PREFlags

var cmdAdminGetPre = &cobra.Command{
	Use:   "pre",
	Short: "Get Preferred Replica Election Details",
	Run: func(cmd *cobra.Command, args []string) {
		var preMeta kafka.PRETopicMeta
		switch true {
		default:
			preMeta = kafka.GetPREMeta(args...)
		}
		switch true {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(preMeta.CreatePREList(), outFmt))
		default:
			kafka.PrintAdm(preMeta.CreatePRESummary())
		}

	},
}

func init() {
	//cmdAdminGetPre.Flags().BoolVar(&preFlags.AllTopics, "alltopics", false, "Filter by Configuration or Key Names.")
}
