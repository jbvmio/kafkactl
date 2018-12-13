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
	"github.com/spf13/cobra"
)

var (
	preAllTopics bool
)

var preCmd = &cobra.Command{
	Use:     "pre",
	Short:   "Perform Preferred Replica Election Tasks",
	Aliases: []string{"preferred-replica-election"},
	PreRun: func(cmd *cobra.Command, args []string) {
		launchZKClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("topic") {
			performTopicPRE(targetTopic)
			return
		}
		if preAllTopics {
			allTopicsPRE()
			return
		}
		closeFatal("Error: Must specify either --topic OR --alltopics, try again.")
	},
}

func init() {
	adminCmd.AddCommand(preCmd)
	preCmd.Flags().BoolVar(&preAllTopics, "alltopics", false, "Initiate Preferred Replica Election for All Topics")
}
