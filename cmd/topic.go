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
	"math/rand"
	"sort"
	"time"

	"github.com/briandowns/spinner"
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
		if meta {
			cmd.Run(metaCmd, args)
			return
		}
		if len(args) == 0 {
			rand.Seed(time.Now().UnixNano())
			r := rand.Intn(43)
			s := spinner.New(spinner.CharSets[r], 100*time.Millisecond)
			t := getAllTopics()
			s.Prefix = fmt.Sprintf("Processing %v Topics ... ", len(t))
			s.Start()
			//sort.Strings(t)
			ts := getTopicSummaries(ctx, t)
			s.Stop()
			sort.Slice(ts, func(i, j int) bool {
				return ts[i].Topic < ts[j].Topic
			})
			formatOutput(ts)
		} else {
			topicList = args
			if targetTopic != "" {
				topicList = append(topicList, targetTopic)
			}
			r := searchForTopics(getAllTopics(), topicList)
			sort.Strings(r)
			if len(r) == 0 {
				fmt.Printf("No results found for topics: %v\n", topicList)
			} else {
				//getTopicMetadata(ctx, r)
				formatOutput(getTopicSummaries(ctx, r))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(topicCmd)
	topicCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
	topicCmd.Flags().BoolVarP(&meta, "metadata", "m", false, "Show metadata details")
}
