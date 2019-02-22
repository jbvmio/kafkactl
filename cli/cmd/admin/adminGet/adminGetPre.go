package adminget

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var preFlags kafka.PREFlags

var cmdAdminGetPre = &cobra.Command{
	Use:     "pre",
	Short:   "Get Preferred Replica Election Details",
	Example: "  kafkactl admin get pre\n  kafkactl admin get pre <topicName>",
	Run: func(cmd *cobra.Command, args []string) {
		var preMeta kafka.PRETopicMeta
		switch {
		default:
			preMeta = kafka.GetPREMeta(args...)
		}
		switch {
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
}
