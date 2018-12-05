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
	"fmt"
	"log"
)

type RAPartList struct {
	Version    int           `json:"version"`
	Partitions []RAPartition `json:"partitions"`
}

type RAPartition struct {
	Topic     string  `json:"topic"`
	Partition int32   `json:"partition"`
	Replicas  []int32 `json:"replicas"`
}

type BrokerReplicas struct {
	BrokerID     int32
	ReplicaCount int32
}

func changeTopicRF(topic string, rFactor int16) {
	exact = true
	tMeta := searchTopicMeta(topic)
	if len(tMeta) > 1 {
		log.Fatalf("No results found for topic: %v\n", topic)
	}
	brokers, err := client.GetClusterMeta()
	if err != nil {
		log.Fatalf("Error retrieving metadata: %v\n", err)
	}
	fmt.Println("BrokerIDs:", brokers.BrokerIDs)
	for _, tm := range tMeta {
		fmt.Println(tm.Topic)
		fmt.Println(tm.Replicas)
	}

	var BR []BrokerReplicas
	available := make(map[int32]bool)
	for _, b := range brokers.BrokerIDs {
		var used bool
		br := BrokerReplicas{
			BrokerID: b,
		}
		for _, tm := range tMeta {
			for _, r := range tm.Replicas {
				if r == b {
					used = true
					br.ReplicaCount++
				}
			}
		}
		if !used {
			available[b] = true
		} else {
			BR = append(BR, br)
		}
	}
	fmt.Println()
	fmt.Printf("BRs: %+v\n", BR)
	fmt.Println("Available:", available)
}
