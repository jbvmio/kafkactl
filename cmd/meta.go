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
	"strings"

	"github.com/spf13/cobra"
)

var metaCmd = &cobra.Command{
	Use:     "meta",
	Short:   "Return Metadata",
	Aliases: []string{"metadata"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && targetTopic == "" {
			cm := getClusterMeta()
			gr := getGroupList()
			c := getController(bootStrap)
			fmt.Println("\nBrokers: ", len(cm.Brokers))
			fmt.Println(" Topics: ", len(cm.Topics))
			fmt.Println(" Groups: ", len(gr))
			fmt.Println("\nCluster:")
			for _, b := range cm.Brokers {
				if strings.Contains(b, c) {
					fmt.Println("*", b)
				} else {
					fmt.Println(" ", b)
				}
			}
			fmt.Printf("\nProvide a topic name to retrieve metadata details for that topic.\n(*)Controller\n")
		} else {
			topicList = args
			if targetTopic != "" {
				topicList = append(topicList, targetTopic)
			}
			r := searchForTopics(getAllTopics(), topicList)
			if len(r) == 0 {
				fmt.Printf("No results found for topics: %v\n", topicList)
			} else {
				getTopicMetadata(ctx, r)
			}
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(metaCmd)
	metaCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
}
