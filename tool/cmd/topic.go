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
	"fmt"

	"github.com/spf13/cobra"
)

var (
	topicList []string
	meta      bool
)

var topicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Search and Retrieve Available Topics",
	Long: `Provides a summary view of available topics.
  Example: kafkactl topic topic1 topic2 topic3
  
If no arguments are provided, all topics are retrieved.
To see detailed metadata information, use the meta command or the -m flag here.
  Example: kafkactl --broker kafkahost --topic topic1 --exact --metadata`,
	Aliases: []string{"topics"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TOPIC")
	},
}

func init() {
	rootCmd.AddCommand(topicCmd)
	topicCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
	topicCmd.Flags().BoolVarP(&meta, "metadata", "m", false, "Show metadata details")
}
