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
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	tBrokers []int32
)

var moveCmd = &cobra.Command{
	Use:     "move",
	Short:   "Move Topic Partitions",
	Aliases: []string{"incr"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("zookeeper") {
			zkServers = []string{zkTargetServer}
		} else {
			if cmd.Flags().Changed("broker") {
				if fileExists(configLocation) {
					_, _, zkServers = tryByBroker(bootStrap, configLocation)
				}
				if len(zkServers) < 1 {
					closeFatal("Error: Increasing Replicas requires access to the corresponding Zookeeper cluster\n  Use --zookeeper zkhost:2181 or configure your ~/.kafkactl config.")
				}
			}
		}
		launchZKClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("topic") {
			if !cmd.Flags().Changed("brokers") || !cmd.Flags().Changed("partitions") {
				closeFatal("Specify both --brokers and --partitions, try again.\n")
			}
			parts := strings.Split(strParts, ",")
			for _, p := range parts {
				tParts = append(tParts, cast.ToInt32(p))
			}
			tpMeta := validateTopicParts(targetTopic, tParts)
			brokes := strings.Split(targetKey, ",")
			for _, b := range brokes {
				tBrokers = append(tBrokers, cast.ToInt32(b))
			}
			validateBrokers(tBrokers)
			movePartitions(tpMeta, tBrokers)
			return
		}
		closeFatal("Must specify --topic, try again.\n")
		return
	},
}

func init() {
	adminCmd.AddCommand(moveCmd)
	moveCmd.Flags().StringVarP(&strParts, "partitions", "p", "", `Comma Separated (eg: "0,1,7,9") Partitions to Move`)
	moveCmd.Flags().StringVar(&targetKey, "brokers", "", `Comma Separated (eg: "1,2,3") brokers to assign replicas / move to`)
	moveCmd.Flags().StringVarP(&zkTargetServer, "zookeeper", "z", "", "Specify a targeted Zookeeper Server and Port (eg. localhost:2181")
}
