// Copyright © 2018 NAME HERE <jbonds@jbvm.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	setConf         bool
	targetConfName  string
	targetConfValue string
)

var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "Get and Set Available Kafka Related Configs (Topics Only for Now)",
	Long:  `Example kafkactl config -t myTopic`,
	Run: func(cmd *cobra.Command, args []string) {
		if targetTopic == "" {
			log.Fatalf("specify a topic, eg. --topic")
		}
		if setConf {
			setTopicConfig(targetTopic, targetConfName, targetConfValue)
			return
		}
		printOutput(searchTopicConfig(targetTopic, targetConfName))
		return

	},
}

func init() {
	adminCmd.AddCommand(confCmd)
	confCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
	confCmd.Flags().BoolVar(&setConf, "set", false, "Alter an Existing Topic Config")
	confCmd.Flags().StringVarP(&targetConfName, "key", "k", "", "Config Option or Key to Set")
	confCmd.Flags().StringVar(&targetConfValue, "value", "", "Config Value to Set")
}
