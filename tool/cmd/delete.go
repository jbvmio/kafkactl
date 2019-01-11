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
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Topics or Groups",
	Long: `Examples:
  kafkactl admin delete -t myTopic
  kafkactl admin delete -g myGroup`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("topic") && !(cmd.Flags().Changed("group") || cmd.Flags().Changed("partition")) {
			closeFatal("specify either group, topic or partition to delete, try again.")
		}
		if cmd.Flags().Changed("topic") {
			if cmd.Flags().Changed("partition") {
				if !cmd.Flags().Changed("offset") {
					closeFatal("must specify the offset to delete to.")
				}
				deleteToOffset(targetTopic, targetPartition, targetOffset)
				return
			}
			deleteTopic(targetTopic)
			return
		}
		if cmd.Flags().Changed("group") {
			deleteGroup(targetGroup)
			return
		}
		closeFatal("specify either --group or --topic, try again.")
	},
}

func init() {
	adminCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().Int32VarP(&targetPartition, "partition", "p", -2, "Partition to Target")
	deleteCmd.Flags().Int64VarP(&targetOffset, "offset", "o", -2, "Desired Offset to Target")
}
