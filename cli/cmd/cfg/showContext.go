package cfg

import (
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var list bool

var cmdShowContext = &cobra.Command{
	Use:     "get-context",
	Aliases: []string{"current-context", "get-contexts"},
	Short:   "Display current and available context details",
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case cmd.CalledAs() == "current-context":
			out.Marshal(GetContext(), outFlags.Format)
		case len(args) < 1:
			out.Marshal(GetContextList(), outFlags.Format)
		default:
			out.Marshal(GetContext(args...), outFlags.Format)
		}
	},
}

func init() {
	//cmdShowContext.Flags().BoolVarP(&list, "list", "l", false, "List available contexts.")
}
