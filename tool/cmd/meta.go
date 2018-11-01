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

var metaCmd = &cobra.Command{
	Use:     "meta",
	Short:   "Return Metadata",
	Aliases: []string{"metadata"},
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := kafkactl.NewClient("aus-glb-kaf03.nix.corp.pps.io:9092")
		defer func() {
			if err := client.Close(); err != nil {
				log.Fatalf("Error closing client: %v\n", err)
			}
		}()
		//client.Logger("")

		/*
			//meta, err := client.GetClusterMeta()
			meta, err := client.GetTopicMeta()
			if err != nil {
				log.Fatalf("Error getting metadata: %v\n", err)
			}
			fmt.Println(meta)
		*/

		//err = client.AddTopic("NewTopicHere", 3, 1)
		//err = client.RemoveTopic("NewTopicHere")
		//err = client.AddPartitions("NewTopicHere", 1)
		//err = client.SetTopicConfig("NewTopicHere", "delete.retention.ms", "43200000")
		//err = client.DeleteToOffset("testtopic", 3, 561)
		err := client.RemoveGroup("jblap")

		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}

		//c, err := client.GetTopicConfig("replicated.test.topic", "flush.messages", "delete.retention.ms")
		//c, err := client.GetTopicConfig("replicated.test.topic")
		//c, err := client.ListTopics()
		c, err := client.ListGroups()
		//c, err := client.BrokerGroups(6)

		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		for _, i := range c {
			fmt.Printf("%+v\n", i)
		}

	},
}

func init() {
	rootCmd.AddCommand(metaCmd)
	metaCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
}
