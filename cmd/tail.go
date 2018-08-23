// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

var (
	targetPartition int
	offset          int64
)

// tailCmd represents the tail command
var tailCmd = &cobra.Command{
	Use:     "tail",
	Short:   "tail or follow output for a specified topic",
	Aliases: []string{"tails", "follow"},
	Run: func(cmd *cobra.Command, args []string) {
		if targetTopic == "" {
			cmd.Help()
			log.Fatalf("\nProvide a topic name to tail it.\n")
		}
		dialer := &kafka.Dialer{
			ClientID:  "kafkaGo",
			DualStack: true,
		}
		p, err := dialer.LookupPartitions(ctx, "tcp", bootStrap, targetTopic)
		if err != nil {
			log.Fatalf("Cannot establish a connection to %v\n", bootStrap)
		}
		var exists bool
		for _, x := range p {
			if targetPartition != -1 {
				if x.ID == targetPartition {
					wg.Add(1)
					go tailPartition(x, dialer)
					exists = true
				}
			} else {
				wg.Add(1)
				go tailPartition(x, dialer)
				exists = true
			}
		}
		if exists != true {
			log.Fatalf("Partition [%v] not found for [%v]\n", targetPartition, targetTopic)
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(tailCmd)

	tailCmd.Flags().IntVarP(&targetPartition, "partition", "p", -1, "Partition to read from")
	tailCmd.Flags().Int64VarP(&offset, "offset", "o", 1, "Relative offset to begin tail")
}
