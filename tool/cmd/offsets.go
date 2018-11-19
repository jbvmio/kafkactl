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

type GroupTopicOffsetMeta struct {
	TopicPartitionMeta
	GroupPartitionMeta
}

type GroupPartitionMeta struct {
	Group            string
	GroupOffset      int64
	Lag              int64
	GroupCoordinator string
}

type TopicPartitionMeta struct {
	Topic       string
	Partition   int32
	TopicOffset int64
	TopicLeader int32
}

func getGroupTopicOffsets(group, topic string) []GroupTopicOffsetMeta {
	exact = true
	var GTMeta []GroupTopicOffsetMeta
	tSum := kafkactl.GetTopicSummaries(searchTopicMeta(topic))
	gList := searchGroupListMeta(group)
	if len(gList) < 1 || len(tSum) < 1 {
		if len(gList) < 1 && len(tSum) >= 1 {
			for _, ts := range tSum {
				for _, p := range ts.Partitions {
					off, err := client.GetOffsetNewest(ts.Topic, p)
					if err != nil {
						off = -7777
					}
					tpm := TopicPartitionMeta{
						Topic:       ts.Topic,
						Partition:   p,
						TopicOffset: off,
					}
					gtMeta := GroupTopicOffsetMeta{
						tpm,
						GroupPartitionMeta{
							Group:            group,
							GroupOffset:      -1001,
							GroupCoordinator: "GRP_UNACTIVE_FOR_TOPIC",
						},
					}
					GTMeta = append(GTMeta, gtMeta)
				}
			}
		} else {
			log.Fatalf("Error retrieving details for specified group and/or topic: [%v > %v]\n", targetGroup, targetTopic)
		}
		return GTMeta
	}
	for _, ts := range tSum {
		for _, p := range ts.Partitions {
			for _, grp := range gList {
				tpm := TopicPartitionMeta{
					Topic:     ts.Topic,
					Partition: p,
				}
				gpm := GroupPartitionMeta{
					Group:            grp.Group,
					GroupCoordinator: grp.Coordinator,
				}
				offset, lag, err := client.OffSetAdmin().Group(grp.Group).Topic(ts.Topic).GetOffsetLag(p)
				if err != nil {
					log.Fatalf("Error: %v\n", err)
				}
				tpm.TopicOffset = (offset + lag)
				gpm.GroupOffset = offset
				gpm.Lag = lag
				gtMeta := GroupTopicOffsetMeta{
					tpm,
					gpm,
				}
				GTMeta = append(GTMeta, gtMeta)
			}
		}
	}
	return GTMeta
}
