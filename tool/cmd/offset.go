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

	"github.com/jbvmio/kafkactl"
	"github.com/spf13/cobra"
)

var (
	targetPartition int32
	targetOffset    int64
)

var offsetCmd = &cobra.Command{
	Use:     "offset",
	Short:   "Get and Set Available Kafka Related Configs (Topics Only for Now)",
	Long:    `Example kafkactl config -t myTopic`,
	Aliases: []string{"offs, os"},
	Run: func(cmd *cobra.Command, args []string) {
		if targetTopic == "" || targetGroup == "" {
			log.Fatalf("specify a group and topic, eg. --topic, --group")
		}
		client, err := kafkactl.NewClient(bootStrap)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		defer func() {
			if err := client.Close(); err != nil {
				log.Fatalf("Error closing client: %v\n", err)
			}
		}()
		if verbose {
			client.Logger("")
		}
		err = client.OffSetAdmin().Group(targetGroup).Topic(targetTopic).ResetOffset(targetPartition, targetOffset)
		if err != nil {
			log.Fatalf("Error reseting offset: %v\n", err)
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(offsetCmd)
	offsetCmd.Flags().Int32VarP(&targetPartition, "partition", "p", -2, "Partition to Target")
	offsetCmd.Flags().Int64VarP(&targetOffset, "offset", "o", -2, "Offset to Target")

	offsetCmd.MarkFlagRequired("partition")
	offsetCmd.MarkFlagRequired("offset")
}
