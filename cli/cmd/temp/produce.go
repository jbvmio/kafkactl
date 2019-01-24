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
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jbvmio/kafkactl"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	targetKey string
	targetMsg string
)

var produceCmd = &cobra.Command{
	Use:     "produce",
	Aliases: []string{"pro", "send"},
	Short:   "Produce messages to a Kafka topic",
	Long: `Produces messages to a Kafka topic.
  Either by:
    Command arguments,
    Stdin,
    Console Producer (launched with no arguments or stdin.)`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("topic") {
			closeFatal("must specify --topic")
		}
		if cmd.Flags().Changed("partitions") || cmd.Flags().Changed("allparts") {
			if cmd.Flags().Changed("partitions") && cmd.Flags().Changed("allparts") {
				closeFatal("Error: Choose either --partitions or --allparts, but not both. Try again.")
			}
			if allParts {
				exact = true
				t := kafkactl.GetTopicSummaries(searchTopicMeta(targetTopic))
				if len(t) != 1 {
					closeFatal("Error: Could not isolate specified topic, %v\n", targetTopic)
				}
				tParts = t[0].Partitions
			}
			if cmd.Flags().Changed("partitions") {
				parts := strings.Split(strParts, ",")
				for _, p := range parts {
					tParts = append(tParts, cast.ToInt32(p))
				}
				validateParts(tParts)
			}
			if cmd.Flags().Changed("message") {
				createMsgSendParts(targetTopic, targetKey, targetMsg, tParts...)
				return
			}
			if stdinAvailable() {
				b, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					closeFatal("Failed to read from stdin: %v\n", err)
				}
				targetMsg = string(bytes.TrimSpace(b))
				createMsgSendParts(targetTopic, targetKey, targetMsg, tParts...)
				return
			}
			launchConsoleProducer(targetTopic, targetKey, tParts...)
			return
		}
		if cmd.Flags().Changed("message") {
			createMsgSend(targetTopic, targetKey, targetMsg, -1)
			return
		}
		if stdinAvailable() {
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				closeFatal("Failed to read from stdin: %v\n", err)
			}
			targetMsg = string(bytes.TrimSpace(b))
			createMsgSend(targetTopic, targetKey, targetMsg, -1)
			return
		}
		launchConsoleProducer(targetTopic, targetKey)
		return
	},
}

func init() {
	rootCmd.AddCommand(produceCmd)
	produceCmd.Flags().StringVarP(&strParts, "partitions", "p", "", `Send to each Comma Separated (eg: "0,1,7,9") Partition)`)
	produceCmd.Flags().BoolVarP(&allParts, "allparts", "a", false, "Send to all Topic Partitions")
	produceCmd.Flags().StringVarP(&targetKey, "key", "k", "", "Optional key to produce/send")
	produceCmd.Flags().StringVarP(&targetMsg, "message", "m", "", "Message/Value to produce/send, Wins over Stdin.")
}
