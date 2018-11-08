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
	resetPartitionOffsetTo(group, topic, partition, int64(-2))
}

func resetPartitionOffsetToNewest(group, topic string, partition int32) {
	resetPartitionOffsetTo(group, topic, partition, int64(-1))
}

func resetPartitionOffsetTo(group, topic string, partition int32, toOffset int64) {
	validPartCheck()
	client, err := kafkactl.NewClient(bootStrap)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("Error closing client: %v\n", err)
		}
	}()
	if verbose {
		client.Logger("")
	}
	err = client.OffSetAdmin().Group(group).Topic(topic).ResetOffset(partition, toOffset)
	if err != nil {
		log.Fatalf("Error reseting offset: %v\n", err)
	}
}

func resetAllPartitionsTo(group, topic string, partition int32, toOffset int64) {
	client, err := kafkactl.NewClient(bootStrap)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("Error closing client: %v\n", err)
		}
	}()
	if verbose {
		client.Logger("")
	}
	exact = true
	ts := kafkactl.GetTopicSummary(searchTopicMeta(targetTopic))
	if len(ts) != 1 {
		log.Fatalf("Error validating topic: %v\n", targetTopic)
	}
	parts := ts[0].Partitions
	if len(parts) < 1 {
		log.Fatalf("No Partitions found for this topic: %v\n", targetTopic)
	}
	var errParts []int32
	for _, p := range parts {
		err = client.OffSetAdmin().Group(group).Topic(topic).ResetOffset(partition, toOffset)
		if err != nil {
			log.Printf("Error reseting offset for partition %v: %v\n", p, err)
			errParts = append(errParts, p)
		}
	}
	if len(errParts) > 1 {
		log.Fatalf("Error reseting offset for the following partitions: %v\n", errParts)
	}
}

func validPartCheck() {
	if targetPartition == -2 {
		log.Fatalf("Error: Specify a valid partition number\n")
	}
}
