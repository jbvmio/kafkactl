package cfg

import (
	"github.com/spf13/cobra"
)

var configFilePath string

var cmdConvertConfig = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"convert-config"},
	Short:   "Convert an older kafkactl config to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		testReplaceOldConfig(configFilePath)
	},
}

func init() {
	cmdConvertConfig.Flags().StringVar(&configFilePath, "filepath", "", "Specify a targeted file, otherwise will look for your ~/.kafkactl.yaml file.")
}
