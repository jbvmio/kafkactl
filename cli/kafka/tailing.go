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

type tailDetails struct {
	topic   string
	partMap map[int32]int64
}

func TailTopic(flags MSGFlags, topics ...string) {
	exact = true
	var count int
	var details []tailDetails
	for _, topic := range topics {
		var parts []int32
		topicSummary := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
		match := true
		switch match {
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
		//endMap := make(map[int32]int64, len(parts))
		offset := getTailValue(flags.Tail)
		for _, p := range parts {
			off, err := client.GetOffsetNewest(topic, p)
			if err != nil {
				closeFatal("Error validating Partition: %v for topic: %v\n", p, err)
			}
			startMap[p] = off + offset
			//endMap[p] = off
			count++
		}
		d := tailDetails{
			topic:   topic,
			partMap: startMap,
		}
		details = append(details, d)
	}
	msgChan := make(chan *kafkactl.Message, 100)
	stopChan := make(chan bool, count)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	for _, d := range details {
		for part, offset := range d.partMap {
			go client.ChanPartitionConsume(d.topic, part, offset, msgChan, stopChan)
		}
	}
ConsumeLoop:
	for {
		select {
		case msg := <-msgChan:
			PrintMSG(msg)
		case <-sigChan:
			fmt.Printf("signal: interrupt\n  Stopping kafkactl ...\n")
			for i := 0; i < count; i++ {
				stopChan <- true
			}
			break ConsumeLoop
		}
	}

	/*
		msgChan := make(chan *kafkactl.Message, 100)
		doneChan := make(chan bool, len(parts))
		for _, p := range parts {
			go func(topic string, p int32) {
				off := startMap[p]
				for off < endMap[p] {
					msg, err := client.ConsumeOffsetMsg(topic, p, off)
					if err != nil {
						closeFatal("Error retrieving message: %v\n", off)
					}
					msgChan <- msg
					off++
				}
				doneChan <- true
			}(topic, p)
		}
		for i := 0; i < len(parts); {
			select {
			case msg := <-msgChan:
				messages = append(messages, msg)
			case <-doneChan:
				i++
			}
		}
	*/

}

/*
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
*/
