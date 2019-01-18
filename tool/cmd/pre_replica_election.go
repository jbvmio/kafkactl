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
	"encoding/json"
	"fmt"
)

type PREList struct {
	Version    int            `json:"version"`
	Partitions []PREPartition `json:"partitions"`
}

type PREPartition struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
}

func performTopicPRE(topic string) {
	exact = true
	tom := chanGetTopicOffsetMap(searchTopicMeta(topic))
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
	zkCreatePRE("/admin/preferred_replica_election", j)
}

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
	//fmt.Printf("%s\n", j)
	zkCreatePRE("/admin/preferred_replica_election", j)
}

func allTopicsPRE() {
	tom := chanGetTopicOffsetMap(searchTopicMeta(""))
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
	zkCreatePRE("/admin/preferred_replica_election", j)
}
