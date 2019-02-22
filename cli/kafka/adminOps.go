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

package kafka

import (
	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/abstraction/kafka"
)

type OpsCreateFlags struct {
	DryRun            bool
	PartitionCount    int32
	ReplicationFactor int16
}

func CreateTopics(partitions int32, rFactor int16, topics ...string) {
	for _, topic := range topics {
		errd = client.AddTopic(topic, partitions, rFactor)
		if errd != nil {
			out.Warnf("Error creating topic: %v", errd)
		} else {
			out.Infof("Successfully created topic %v", topic)
		}
	}
}

func DeleteTopics(topics ...string) {
	for _, topic := range topics {
		errd = client.RemoveTopic(topic)
		if errd != nil {
			out.Warnf("Error deleting topic: %v", errd)
		} else {
			out.Infof("Successfully deleted topic %v", topic)
		}
	}
}

func DeleteGroups(groups ...string) {
	for _, group := range groups {
		errd = client.RemoveGroup(group)
		if errd != nil {
			out.Warnf("Error deleting group: %v", errd)
		} else {
			out.Infof("Successfully deleted group %v", group)
		}
	}
}

func ConfigurePartitionCount(flags OpsCreateFlags, topics ...string) {
	if !ClientVersion().IsAtLeast(kafkactl.MinCreatePartsVer) {
		closeFatal("Error: Required Version %v needed for create partitions. Using: %v", kafkactl.MinCreatePartsVer, ClientVersion())
	}
	for _, topic := range topics {
		err := client.AddPartitions(topic, flags.PartitionCount, flags.DryRun)
		switch {
		case flags.DryRun:
			switch {
			case err != nil:
				out.Warnf("Dry Run:\nERROR Adding Partitions: %v", err)
			default:
				out.Infof("Dry Run:\nIncrease Partitions Successful - Topic: %v Partitions: %v", topic, flags.PartitionCount)
			}
		case err != nil:
			out.Warnf("ERROR Adding Partitions: %v", err)
		default:
			out.Infof("Increase Partitions Successful - Topic: %v Partitions: %v", topic, flags.PartitionCount)
		}
	}
}
