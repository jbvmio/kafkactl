package lag

import (
	"strings"

	"github.com/jbvmio/kafkactl/cli/kafka"
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var lagFlags kafka.LagFlags
var onlyTotal bool

var ValidOffsets bool

var CmdGetLag = &cobra.Command{
	Use:     "lag",
	Short:   "Get Lag Info",
	Example: examples.GetLag(),
	Run: func(cmd *cobra.Command, args []string) {
		var lag []kafka.PartitionLag
		var totalLag []kafka.TotalLag
		switch true {
		case strings.Contains(cmd.CalledAs(), "topic"):
			lag = kafka.GetGroupLag(kafka.GroupMetaByTopics(args...))
		case cmd.CalledAs() == "member":
			lag = kafka.GetGroupLag(kafka.GroupMetaByMember(args...))
		case len(args) > 0:
			lag = kafka.GetGroupLag(kafka.SearchGroupMeta(args...))
		default:
			totalLag = kafka.FindTotalLag()
			onlyTotal = true
		}
		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			if onlyTotal {
				out.IfErrf(out.Marshal(totalLag, outFmt))
				return
			}
			out.IfErrf(out.Marshal(lag, outFmt))
		default:
			if onlyTotal {
				kafka.PrintOut(totalLag)
				return
			}
			kafka.PrintOut(lag)
		}
	},
}

func init() {
}
