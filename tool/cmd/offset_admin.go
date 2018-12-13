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
)

func resetPartitionOffsetToOldest(group, topic string, partition int32) {
	resetPartitionOffsetTo(group, topic, partition, kafkactl.OffsetOldest)
}

func resetPartitionOffsetToNewest(group, topic string, partition int32) {
	resetPartitionOffsetTo(group, topic, partition, kafkactl.OffsetNewest)
}

func resetPartitionOffsetTo(group, topic string, partition int32, toOffset int64) {
	validPartCheck()
	if toOffset == kafkactl.OffsetNewest {
		toOffset, errd = client.GetOffsetNewest(topic, partition)
		if errd != nil {
			closeFatal("Error retrieving newest offset: %v\n", errd)
		}
	}
	if toOffset == kafkactl.OffsetOldest {
		toOffset, errd = client.GetOffsetOldest(topic, partition)
		if errd != nil {
			closeFatal("Error retrieving oldest offset: %v\n", errd)
		}
	}
	errd = client.OffSetAdmin().Group(group).Topic(topic).ResetOffset(partition, toOffset)
	if errd != nil {
		closeFatal("Error reseting offset: %v\n", errd)
	}
}

func resetAllPartitionsTo(group, topic string, toOffset int64) {
	exact = true
	ts := kafkactl.GetTopicSummaries(searchTopicMeta(targetTopic))
	if len(ts) != 1 {
		closeFatal("Error validating topic: %v\n", targetTopic)
	}
	parts := ts[0].Partitions
	if len(parts) < 1 {
		closeFatal("No Partitions found for this topic: %v\n", targetTopic)
	}
	var errParts []int32
	for _, p := range parts {
		if toOffset == kafkactl.OffsetNewest {
			toOffset, errd = client.GetOffsetNewest(topic, p)
			if errd != nil {
				closeFatal("Error retrieving newest offset: %v\n", errd)
			}
		}
		if toOffset == kafkactl.OffsetOldest {
			toOffset, errd = client.GetOffsetOldest(topic, p)
			if errd != nil {
				closeFatal("Error retrieving oldest offset: %v\n", errd)
			}
		}
		errd = client.OffSetAdmin().Group(group).Topic(topic).ResetOffset(p, toOffset)
		if errd != nil {
			log.Printf("Error reseting offset for partition %v: %v\n", p, errd)
			errParts = append(errParts, p)
		}
	}
	if len(errParts) > 1 {
		closeFatal("Error reseting offset for the following partitions: %v\n", errParts)
	}
}

func validPartCheck() {
	if targetPartition == -2 {
		closeFatal("Error: Specify a valid partition number\n")
	}
}
