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
	"encoding/json"

	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/jbvmio/kafkactl"
	"github.com/jbvmio/kafkactl/cli/zookeeper"
)

type PREFlags struct {
	AllTopics  bool
	Partition  int32
	Partitions []string
}

type PRESummary struct {
	Topics   []string
	PRECount map[string]uint32
}

type PRETopicMeta struct {
	Partitions []kafkactl.TopicMeta
}

type PREList struct {
	Version    int            `json:"version"`
	Partitions []PREPartition `json:"partitions"`
}

type PREPartition struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
}

const (
	prePath = `/admin/preferred_replica_election`
)

func PerformTopicPRE(topics ...string) PRETopicMeta {
	preMeta := GetPREMeta(topics...)
	if len(preMeta.Partitions) < 1 {
		closeFatal("Error: Unable to isolate topics: %v\n", topics)
	}
	j, err := json.Marshal(preMeta.CreatePREList())
	if err != nil {
		closeFatal("Error Marshaling Topic/Partition Data: %v\n", err)
	}
	zkCreatePRE(j)
	return preMeta
	//zkCreatePRE("/admin/preferred_replica_election", j)
}

func zkCreatePRE(data []byte) {
	handleC("%v", zookeeper.KafkaZK(targetContext, verbose))
	check, err := zkClient.Exists(prePath)
	handleC("Error: %v", err)
	if check {
		closeFatal("Preferred Replica Election Already in Progress.")
	}
	zookeeper.ZKCreate(prePath, false, data...)
}

func GetPREMeta(topics ...string) PRETopicMeta {
	var preMeta []kafkactl.TopicMeta
	tMeta := SearchTopicMeta(topics...)
	for _, tm := range tMeta {
		if len(tm.Replicas) > 0 {
			if tm.Leader != tm.Replicas[0] {
				preMeta = append(preMeta, tm)
			}
		}
	}
	if len(preMeta) < 1 {
		CloseClient()
		out.Exitf(0, "No Preferred Replica Elections Needed.")
	}
	return PRETopicMeta{Partitions: preMeta}
}

func (pre PRETopicMeta) CreatePREList() PREList {
	var preParts []PREPartition
	for _, part := range pre.Partitions {
		preP := PREPartition{
			Topic:     part.Topic,
			Partition: part.Partition,
		}
		preParts = append(preParts, preP)
	}
	return PREList{Version: 1, Partitions: preParts}
}

func (pre PRETopicMeta) CreatePRESummary() PRESummary {
	var topicList []string
	topicDone := make(map[string]bool)
	preCount := make(map[string]uint32)
	for _, part := range pre.Partitions {
		preCount[part.Topic]++
		if !topicDone[part.Topic] {
			topicList = append(topicList, part.Topic)
			topicDone[part.Topic] = true
		}
	}
	return PRESummary{Topics: topicList, PRECount: preCount}
}

/*

func preTopicsStdin(td []topicStdinData) {
	var preParts []PREPartition
	for _, t := range td {

		preP := PREPartition{
			Topic:     t.topic,
			Partition: t.partition,
		}
		preParts = append(preParts, preP)

	}
	preList := PREList{
		Version:    1,
		Partitions: preParts,
	}
	j, err := json.Marshal(preList)
	if err != nil {
		closeFatal("Error Marshaling Topic/Partition Data: %v\n", err)
	}
	fmt.Printf("%s", j)
	//zkCreatePRE("/admin/preferred_replica_election", j)
}

func performTopicPRE(topic string) {
	exact = true
	tom := GetTopicOffsetMap(SearchTopicMeta(topic))
	if len(tom) < 1 {
		closeFatal("Error: Cannot find topic: %v\n", topic)
	}
	var preParts []PREPartition
	for _, to := range tom {
		for k := range to.PartitionOffsets {
			preP := PREPartition{
				Topic:     to.Topic,
				Partition: k,
			}
			preParts = append(preParts, preP)
		}
	}
	preList := PREList{
		Version:    1,
		Partitions: preParts,
	}
	j, err := json.Marshal(preList)
	if err != nil {
		closeFatal("Error Marshaling Topic/Partition Data: %v\n", err)
	}
	fmt.Printf("%s", j)
	//zkCreatePRE("/admin/preferred_replica_election", j)
}

func allTopicsPRE() {
	tom := GetTopicOffsetMap(SearchTopicMeta())
	if len(tom) < 1 {
		closeFatal("Error: No Topics Available.\n")
	}
	var preParts []PREPartition
	for _, to := range tom {
		for k := range to.PartitionOffsets {
			preP := PREPartition{
				Topic:     to.Topic,
				Partition: k,
			}
			preParts = append(preParts, preP)
		}
	}
	preList := PREList{
		Version:    1,
		Partitions: preParts,
	}
	j, err := json.Marshal(preList)
	if err != nil {
		closeFatal("Error Marshaling Topic/Partition Data: %v\n", err)
	}
	fmt.Printf("%s", j)
	//zkCreatePRE("/admin/preferred_replica_election", j)
}

*/
