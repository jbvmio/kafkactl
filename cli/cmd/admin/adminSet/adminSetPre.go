package adminset

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var preFlags kafka.PREFlags

var cmdAdminSetPre = &cobra.Command{
	Use:     "pre",
	Example: "  kafkactl admin set pre\n  kafkactl admin set pre <topicName>",
	Short:   "Set Preferred Replica Elections",
	Run: func(cmd *cobra.Command, args []string) {
		var preMeta kafka.PRETopicMeta
		switch true {
		default:
			preMeta = kafka.PerformTopicPRE(args...)
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
}
