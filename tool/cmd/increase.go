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

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var increaseCmd = &cobra.Command{
	Use:     "increase",
	Short:   "Increase Topic Partitions / Replication Factor",
	Aliases: []string{"incr"},
	PreRun: func(cmd *cobra.Command, args []string) {
		launchZKClient()
		nonMainCMD = true
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("topic") {
			if !cmd.Flags().Changed("replicas") {
				log.Fatalf("Use --replicas to specify the total number of replicas desired for topic: %v\n", targetTopic)
			}
			rFactor := cast.ToInt(targetRFactor)
			performPartitionReAssignment(targetTopic, rFactor)
			return
		}
		return
	},
}

func init() {
	adminCmd.AddCommand(increaseCmd)
	increaseCmd.Flags().Int16VarP(&targetRFactor, "replicas", "r", -2, "Desired, Total Number of Replicas")
	//increaseCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
}
