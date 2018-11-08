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

	"github.com/spf13/cobra"
)

var (
	targetOffset int64
	performReset bool
	allParts     bool
	targetNewest bool
	targetOldest bool
)

var offsetCmd = &cobra.Command{
	Use:     "offset",
	Short:   "Manage Offsets",
	Long:    `Example kafkactl admin offset -g myGroup -t myTopic -p 5 --reset-to 999`,
	Aliases: []string{"offsets"},
	Run: func(cmd *cobra.Command, args []string) {
		if targetTopic == "" || targetGroup == "" {
			log.Fatalf("specify both group and topic")
		}
		gto, err := getGroupTopicOffsets(targetGroup, targetTopic)
		if err != nil {
			log.Fatalf("ERROR: %v\n", err)
		}
		fmt.Printf("%+v\n", gto)
		/*
			exact = true
			showLag = true
			desc := []string{"group"}
			desc = append(desc, targetGroup)
			describeCmd.Run(cmd, desc)
			return
		*/
	},
}

func init() {
	adminCmd.AddCommand(offsetCmd)
	offsetCmd.Flags().Int32VarP(&targetPartition, "partition", "p", -2, "Partition to Target")
	offsetCmd.Flags().Int64VarP(&targetOffset, "offset", "o", -2, "Offset to Target")
	offsetCmd.Flags().BoolVar(&performReset, "reset", false, "Initiate a ResetOffset Operation")
	offsetCmd.Flags().BoolVar(&targetNewest, "newest", false, "Target the Next Msg that would be Produced to the Topic")
	offsetCmd.Flags().BoolVar(&targetOldest, "oldest", false, "Target the Oldest Msg still available on the Topic")

	//offsetCmd.MarkFlagRequired("partition")
	//offsetCmd.MarkFlagRequired("offset")
}
