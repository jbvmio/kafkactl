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
	"github.com/jbvmio/kafkactl"
)

// PartitionLag struct def:
type PartitionLag struct {
	Group     string
	Topic     string
	Partition int32
	Member    string
	Offset    int64
	Lag       int64
}

func getPartitionLag(grpMeta []kafkactl.GroupMeta) []PartitionLag {
	var partitionLag []PartitionLag
	for _, gm := range grpMeta {
		for _, m := range gm.MemberAssignments {
			for topic, partitions := range m.TopicPartitions {
				for _, p := range partitions {
					offset, lag, err := client.OffSetAdmin().Group(gm.Group).Topic(topic).GetOffsetLag(p)
					if err != nil {
						lag = -7777
					}
					pl := PartitionLag{
						Group:     gm.Group,
						Topic:     topic,
						Partition: p,
						Member:    m.ClientID,
						Offset:    offset,
						Lag:       lag,
					}
					partitionLag = append(partitionLag, pl)
				}
			}
		}
	}
	return partitionLag
}

func chanGetPartitionLag(grpMeta []kafkactl.GroupMeta) []PartitionLag {
	var partitionLag []PartitionLag
	for _, gm := range grpMeta {
		for _, m := range gm.MemberAssignments {
			for topic, partitions := range m.TopicPartitions {
				plChan := make(chan PartitionLag, 100)
				for _, p := range partitions {
					go goGetPartitionLag(client, gm.Group, topic, m.ClientID, p, plChan)
				}
				for i := 0; i < len(partitions); i++ {
					select {
					case pl := <-plChan:
						partitionLag = append(partitionLag, pl)
					}
				}
			}
		}
	}
	return partitionLag
}

func goGetPartitionLag(client *kafkactl.KClient, group, topic, clientID string, partition int32, plChan chan PartitionLag) {
	offset, lag, err := client.OffSetAdmin().Group(group).Topic(topic).GetOffsetLag(partition)
	if err != nil {
		lag = -7777
	}
	pl := PartitionLag{
		Group:     group,
		Topic:     topic,
		Partition: partition,
		Member:    clientID,
		Offset:    offset,
		Lag:       lag,
	}
	plChan <- pl
}
