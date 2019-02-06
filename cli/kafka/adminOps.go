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
)

type OpsCreateFlags struct {
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

func DeleteGroup(group string) {
	errd = client.RemoveGroup(group)
	handleC("Error removing group: %v", errd)
	out.Infof("\nSuccessfully removed group %v\n", group)
}

func DeleteToOffset(topic string, partition int32, offset int64) {
	errd = client.DeleteToOffset(topic, partition, offset)
	handleC("Error deleting to offset: %v", errd)
	out.Infof("\nSuccessfully deleted to offset %v", offset)
}
