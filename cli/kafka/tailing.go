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
	"time"

	"github.com/Shopify/sarama"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/tidwall/gjson"

	kafkactl "github.com/jbvmio/kafka"
)

type followDetails struct {
	topic   string
	partMap map[int32]int64
}

func FollowTopic(flags MSGFlags, outFlags out.OutFlags, topics ...string) {
	exact = true
	var count int
	var details []followDetails
	var timeCheck time.Time
	var useFilter bool
	if len(flags.JSONFilters) > 0 {
		useFilter = true
	}
	for _, topic := range topics {
		var parts []int32
		topicSummary := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
		switch true {
		case len(topicSummary) != 1:
			closeFatal("Error isolating topic: %v\n", topic)
		case flags.Partition != -1:
			parts = append(parts, flags.Partition)
		case len(flags.Partitions) == 0:
			parts = topicSummary[0].Partitions
		default:
			parts = validateParts(flags.Partitions)
		}
		startMap := make(map[int32]int64, len(parts))
		offset := getTailValue(flags.Tail)
		for _, p := range parts {
			off, err := client.GetOffsetNewest(topic, p)
			if err != nil {
				closeFatal("Error validating Partition: %v for topic: %v\n", p, err)
			}
			startMap[p] = off + offset
			count++
		}
		d := followDetails{
			topic:   topic,
			partMap: startMap,
		}
		details = append(details, d)
	}
	msgChan := make(chan *sarama.ConsumerMessage, 100)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	for _, d := range details {
		for part, offset := range d.partMap {
			go client.ChanPartitionConsume(d.topic, part, offset, msgChan)
		}
	}
ConsumeLoop:
	for {
		select {
		case <-sigChan:
			fmt.Printf("signal: interrupt\n  Stopping kafkactl ... ")
			break ConsumeLoop
		case msg := <-msgChan:
			switch {
			case msg.Timestamp == timeCheck:
				if len(msg.Value) != 0 {
					out.Warnf("%s", msg.Value)
				}
			case useFilter:
				fmt.Printf("%s\n\n", gjson.GetManyBytes(msg.Value, flags.JSONFilters...))
			default:
				PrintMSG(msg, outFlags)
			}
			continue ConsumeLoop
		}
	}
	client.StopPartitionConsumers()
	fmt.Println("Stopped")
}
