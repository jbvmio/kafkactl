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
	showMsgKey       bool
	showMsgTimestamp bool
	timeQuery        string
)

var messageCmd = &cobra.Command{
	Use:     "message",
	Aliases: []string{"msg", "read", "getmsg"},
	Short:   "Retreive Targeted Messages from a Kafka Topic",
	Long:    `Example: kafkactl message -t myTopic -p 5 -o 6723`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("topic") || !cmd.Flags().Changed("partition") {
			log.Fatalf("must specify 1 --topic, --partition and --offset OR --timestamp")
		}
		if cmd.Flags().Changed("offset") && cmd.Flags().Changed("timestamp") {
			log.Fatalf("must specify 3 --offset OR --timestamp, not both. Try again.")
		}
		if cmd.Flags().Changed("timestamp") {
			msg := getMSGByTime(targetTopic, targetPartition, timeQuery)
			printOutput(msg)
			return
		}
		msg := getMSG(targetTopic, targetPartition, targetOffset)
		printOutput(msg)
		return
	},
}

func init() {
	rootCmd.AddCommand(messageCmd)
	messageCmd.Flags().Int32VarP(&targetPartition, "partition", "p", 0, "Target Partition to Retrieve Message")
	messageCmd.Flags().Int64VarP(&targetOffset, "offset", "o", 0, "Target Offset to Retrieve Message")
	messageCmd.Flags().StringVarP(&timeQuery, "timestamp", "q", "", `Retrieve by Timestamp (Ex: "11/10/2018 15:30:000")`)
	messageCmd.Flags().BoolVarP(&showMsgKey, "key", "k", false, "Show Msg Key Values")
	messageCmd.Flags().BoolVar(&showMsgTimestamp, "show-time", false, "Show Msg Timestamps")
}
