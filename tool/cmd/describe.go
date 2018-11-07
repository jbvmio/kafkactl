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
	"log"
	"os"
	"strings"

	"github.com/jbvmio/kafkactl"

	"github.com/spf13/cobra"
)

var (
	showLag     bool
	showConfig  bool
	targetMatch = true
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Return Topic or Group details",
	Long: `Examples:
  kafkactl describe group myConsumerGroup
  kafkactl describe topic myTopic`,
	Aliases: []string{"desc", "descr", "des"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}
		target := args[0]
		args = args[1:]
		switch targetMatch {
		case strings.Contains(target, "gro"):
			grps := searchGroupMeta(args...)
			if targetTopic != "" {
				grps = groupMetaByTopic(targetTopic, grps)
			}
			if len(grps) < 1 {
				log.Fatalf("no results for that group/topic combination\n")
			}
			if showLag {
				var partitionLag []PartitionLag
				if useFast {
					partitionLag = chanGetPartitionLag(grps)
				} else {
					partitionLag = getPartitionLag(grps)
				}
				printOutput(partitionLag)
				return
			}
			printOutput(grps)
			return
		case strings.Contains(target, "top"):
			if showConfig {
				var topics []string
				ts := kafkactl.GetTopicSummary(searchTopicMeta(args...))
				for _, t := range ts {
					topics = append(topics, t.Topic)
				}
				configs := getTopicConfig(topics...)
				printOutput(configs)
				return
			}
			tm := searchTopicMeta(args...)
			if len(tm) < 1 {
				log.Fatalf("no results for that group/topic combination\n")
			}
			printOutput(tm)
			return
		default:
			fmt.Println("no such resource to describe:", target)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
	describeCmd.Flags().BoolVarP(&showLag, "lag", "l", false, "Display Offset and Lag details")
	describeCmd.Flags().BoolVar(&showConfig, "conf", false, "Show Configuration details (Topics)")
	//describeCmd.Flags().StringVarP(&clientID, "clientid", "i", "", "Find groups by ClientID")
}
