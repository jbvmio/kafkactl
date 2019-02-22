package cfg

import (
	examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var outFlags out.OutFlags
var showSample bool

var CmdConfig = &cobra.Command{
	Use:     "config",
	Example: examples.Config(),
	Short:   "Show and Edit kafkactl config",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		if showSample {
			GenSample()
		}
	},
}

func init() {
	CmdConfig.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "yaml", "Output Format - yaml|json.")
	CmdConfig.Flags().BoolVar(&showSample, "sample", false, "Display a sample config file.")

	CmdConfig.AddCommand(cmdView)
	CmdConfig.AddCommand(cmdShowContext)
	CmdConfig.AddCommand(cmdUseContext)
	CmdConfig.AddCommand(cmdConvertConfig)
}
