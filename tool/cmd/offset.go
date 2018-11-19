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

	"github.com/jbvmio/kafkactl"

	"github.com/spf13/cobra"
)

var (
	targetOffset   int64
	targetRelative int64
	performReset   bool
	allParts       bool
	targetNewest   bool
	targetOldest   bool
)

var offsetCmd = &cobra.Command{
	Use:   "offset",
	Short: "Manage Offsets",
	Long: `Example: kafkactl admin offset -g myGroup -t myTopic -p 5 --reset-to 999
  Default: (without using --reset flag) shows partition offset details for the specified group and topic.`,
	Aliases: []string{"offsets"},
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("topic") || !cmd.Flags().Changed("group") {
			log.Fatalf("specify both group and topic")
		}
		if performReset {
			if cmd.Flags().Changed("newest") && cmd.Flags().Changed("oldest") {
				log.Fatalf("cannot specify both newest and oldest, try again.")
			}
			if cmd.Flags().Changed("allparts") && cmd.Flags().Changed("partition") {
				log.Fatalf("cannot specify both allParts and single partition, try again.")
			}
			if targetNewest || targetOldest {
				if cmd.Flags().Changed("offset") || cmd.Flags().Changed("reletive") {
					log.Fatalf("cannot specify specific offset using either newest/oldest, try again.")
				}
				if targetNewest {
					if allParts {
						resetAllPartitionsTo(targetGroup, targetTopic, kafkactl.OffsetNewest)
						return
					}
					if !cmd.Flags().Changed("partition") {
						log.Fatalf("need a valild partition, try again.")
					}
					resetPartitionOffsetToNewest(targetGroup, targetTopic, targetPartition)
					return
				}
				if targetOldest {
					if allParts {
						resetAllPartitionsTo(targetGroup, targetTopic, kafkactl.OffsetOldest)
						return
					}
					if !cmd.Flags().Changed("partition") {
						log.Fatalf("need a valild partition, try again.")
					}
					resetPartitionOffsetToOldest(targetGroup, targetTopic, targetPartition)
					return
				}
			}
			if cmd.Flags().Changed("reletive") {
				fmt.Println("RELETIVE WiP*")
				return
			}
			if allParts {
				resetAllPartitionsTo(targetGroup, targetTopic, targetOffset)
				return
			}
			if !cmd.Flags().Changed("offset") || !cmd.Flags().Changed("partition") {
				log.Fatalf("specify both --partition and --offset OR --allparts, try again.")
			}
			resetPartitionOffsetTo(targetGroup, targetTopic, targetPartition, targetOffset)
			return
		}
		gto := getGroupTopicOffsets(targetGroup, targetTopic)
		printOutput(gto)
		return
	},
}

func init() {
	adminCmd.AddCommand(offsetCmd)
	offsetCmd.Flags().Int32VarP(&targetPartition, "partition", "p", -2, "Partition to Target")
	offsetCmd.Flags().Int64VarP(&targetOffset, "offset", "o", -2, "Desired Offset")
	offsetCmd.Flags().Int64Var(&targetRelative, "relative", 0, "Reletive offset from the Topic Offset (Must be Negative, eg. -5)")
	offsetCmd.Flags().BoolVarP(&allParts, "allparts", "a", false, "Perform on All Available Partitions for Specified Group and Topic.")
	offsetCmd.Flags().BoolVar(&performReset, "reset", false, "Initiate a ResetOffset Operation")
	offsetCmd.Flags().BoolVar(&targetNewest, "newest", false, "Target the Next Msg that would be Produced to the Topic")
	offsetCmd.Flags().BoolVar(&targetOldest, "oldest", false, "Target the Oldest Msg still available on the Topic")
}
