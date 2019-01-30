package adminset

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/cobra"
)

var preFlags kafka.PREFlags

var cmdAdminSetPre = &cobra.Command{
	Use:   "pre",
	Short: "Set Preferred Replica Elections",
	Run: func(cmd *cobra.Command, args []string) {
		match := true
		switch match {
		default:
			out.Infof("SET PRE HERE")
		}

		/*
			switch match {
			case cmd.Flags().Changed("out"):
				outFmt, err := cmd.Flags().GetString("out")
				if err != nil {
					out.Warnf("WARN: %v", err)
				}
				out.IfErrf(out.Marshal(topicConfigs, outFmt))
			default:
				kafka.PrintOut(topicConfigs)
			}
		*/
	},
}

func init() {
	//cmdAdminSetPre.Flags().BoolVar(&preFlags.AllTopics, "alltopics", false, "Filter by Configuration or Key Names.")
}
