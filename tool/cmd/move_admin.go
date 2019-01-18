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

	"github.com/jbvmio/kafkactl"
)

func validateBrokers(brokers []int32) {
	pMap := make(map[int32]bool, len(brokers))
	for _, p := range brokers {
		if pMap[p] {
			closeFatal("Error: invalid broker entered or duplicate.")
		} else {
			pMap[p] = true
		}
	}
	cm, err := client.GetClusterMeta()
	if err != nil {
		closeFatal("Error validating brokers: %v\n", err)
	}
	bMap := make(map[int32]bool, len(brokers))
	for _, b := range brokers {
		for _, broker := range cm.BrokerIDs {
			if broker == b {
				bMap[b] = true
			}
		}
	}
	for _, b := range brokers {
		if !bMap[b] {
			closeFatal("Error validating broker: Invalid Broker > %v\n", b)
		}
	}
}

func validateTopicParts(topic string, parts []int32) []kafkactl.TopicMeta {
	pMap := make(map[int32]bool, len(parts))
	for _, p := range parts {
		if pMap[p] {
			closeFatal("Error: invalid partition entered or duplicate.")
		} else {
			pMap[p] = true
		}
	}
	exact = true
	tMeta := searchTopicMeta(topic)
	if len(tMeta) < 1 {
		closeFatal("Error validating topic: %v\n", topic)
	}
	filterTMeta := make([]kafkactl.TopicMeta, 0, len(parts))
	tMap := make(map[int32]bool, len(parts))
	for _, p := range parts {
		for _, tm := range tMeta {
			if tm.Partition == p {
				tMap[p] = true
				filterTMeta = append(filterTMeta, tm)
			}
		}
	}
	for _, p := range parts {
		if !tMap[p] {
			closeFatal("Error validating partition: Invalid Partition -> %v\n", p)
		}
	}
	return filterTMeta
}

func movePartitions(topicMeta []kafkactl.TopicMeta, brokers []int32) {
	var raparts []RAPartition
	for _, tm := range topicMeta {
		rap := RAPartition{
			Topic:     tm.Topic,
			Partition: tm.Partition,
			Replicas:  brokers,
		}
		raparts = append(raparts, rap)
	}
	rapList := RAPartList{
		Version:    1,
		Partitions: raparts,
	}
	j, err := json.Marshal(rapList)
	if err != nil {
		closeFatal("Error on Marshal: %v\n", err)
	}
	zkCreateReassignPartitions("/admin/reassign_partitions", j)
}

func movePartitionsStdin(moveData []topicStdinData, brokers []int32) {
	var raparts []RAPartition
	for _, tm := range moveData {
		rap := RAPartition{
			Topic:     tm.topic,
			Partition: tm.partition,
			Replicas:  brokers,
		}
		raparts = append(raparts, rap)
	}
	rapList := RAPartList{
		Version:    1,
		Partitions: raparts,
	}
	j, err := json.Marshal(rapList)
	if err != nil {
		closeFatal("Error on Marshal: %v\n", err)
	}
	fmt.Printf("%s", j)
	zkCreateReassignPartitions("/admin/reassign_partitions", j)
}
