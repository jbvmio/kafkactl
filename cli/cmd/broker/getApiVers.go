package broker

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var CmdGetApiVers = &cobra.Command{
	Use:     "apis",
	Aliases: []string{"api"},
	Short:   "Get API Protocol Version Details",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case cmd.Flags().Changed("out"):
			outFmt, err := cmd.Flags().GetString("out")
			if err != nil {
				out.Warnf("WARN: %v", err)
			}
			out.IfErrf(out.Marshal(kafka.GetAPIVersions(), outFmt))
		default:
			kafka.PrintOut(kafka.GetAPIVersions())
		}
	},
}

func init() {
}
