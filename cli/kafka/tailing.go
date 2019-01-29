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
	"fmt"
	"os"
	"os/signal"

	"github.com/jbvmio/kafkactl"
)

func getTopicMsg(topic string, partition int32, offset int64) {
	msg, err := client.ConsumeOffsetMsg("testtopic", 0, 1955)
	if err != nil {
		closeFatal("Error: %v\n", err)
	}
	fmt.Printf("%s", msg.Value)
}

func tailTopic(topic string, relativeOffset int64, partitions ...int32) {
	if relativeOffset > 0 {
		closeFatal("reletive offset must be a negative number")
	}
	exact = true
	tSum := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
	if len(tSum) != 1 {
		closeFatal("Error finding topic: %v\n", topic)
	}
	if len(partitions) == 0 {
		partitions = tSum[0].Partitions
	}
	pMap := make(map[int32]int64)
	for _, ts := range tSum {
		for _, p := range partitions {
			off, err := client.GetOffsetNewest(ts.Topic, p)
			if err != nil {
				closeFatal("Error validating Partition: %v for topic: %v\n", p, err)
			}
			pMap[p] = off
		}
	}
	msgChan := make(chan *kafkactl.Message, 100)
	stopChan := make(chan bool, len(pMap))
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	for part, offset := range pMap {
		rel := (offset + relativeOffset)
		go client.ChanPartitionConsume(topic, part, rel, msgChan, stopChan)
	}
ConsumeLoop:
	for {
		select {
		case msg := <-msgChan:
			fmt.Printf("%s\n", msg.Value)
		case <-sigChan:
			fmt.Printf("signal: interrupt\n  Stopping kafkactl ...\n")
			for i := 0; i < len(pMap); i++ {
				stopChan <- true
			}
			break ConsumeLoop
		}
	}
}
