package lag

import (
	"strings"

	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var lagFlags kafka.LagFlags

var CmdGetLag = &cobra.Command{
	Use:   "lag",
	Short: "Get Lag Info",
	Run: func(cmd *cobra.Command, args []string) {
		var lag []kafka.PartitionLag
		switch true {
		case strings.Contains(cmd.CalledAs(), "topic"):
			lag = kafka.GetGroupLag(kafka.GroupMetaByTopics(args...))
		case cmd.CalledAs() == "member":
			lag = kafka.GetGroupLag(kafka.GroupMetaByMember(args...))
		case len(args) > 0:
			lag = kafka.GetGroupLag(kafka.SearchGroupMeta(args...))
		default:
			lag = kafka.FindLag()
		}
		switch true {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(lag, outFmt))
		default:
			kafka.PrintOut(lag)
		}
	},
}

func init() {
}
