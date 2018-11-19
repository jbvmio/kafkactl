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
	targetRFactor int16
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Topics",
	Long:  `Example kafkactl admin create -t myTopic`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("topic") || !cmd.Flags().Changed("partitions") || !cmd.Flags().Changed("rfactor") {
			log.Fatalf("must specify --topic, --partitions and --rfactor")
		}
		createTopic(targetTopic, targetPartition, targetRFactor)
		return
	},
}

func init() {
	adminCmd.AddCommand(createCmd)
	createCmd.Flags().Int32VarP(&targetPartition, "partitions", "p", -2, "Desired Partition Count")
	createCmd.Flags().Int16VarP(&targetRFactor, "rfactor", "r", -2, "Desired Replication Factor")
}
