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
	"github.com/jbvmio/kafkactl/cli/x/out"
)

type MSGFlags struct {
	Partitions  []string
	Partition   int32
	Offset      int64
	Tail        int64
	TailTouched bool
	Follow      bool
}

func GetMessages(flags MSGFlags, topics ...string) []*kafkactl.Message {
	exact = true
	var messages []*kafkactl.Message
	switch true {
	case flags.TailTouched:
		return tailMSGs(flags, topics...)
	default:
		return getMSGs(flags, topics...)
	}
	return messages
}

func getMSGs(flags MSGFlags, topics ...string) []*kafkactl.Message {
	var messages []*kafkactl.Message
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
		pMap := make(map[int32]int64, len(parts))
		switch true {
		case flags.Offset == -1:
			for _, p := range parts {
				off, err := client.GetOffsetNewest(topic, p)
				if err != nil {
					closeFatal("Error validating Partition: %v for topic: %v\n", p, err)
				}
				pMap[p] = off + flags.Offset
			}
		default:
			for _, p := range parts {
				pMap[p] = flags.Offset
			}
		}
		for part, off := range pMap {
			msg, err := client.ConsumeOffsetMsg(topic, part, off)
			if err != nil {
				out.Warnf("WARN %v [%v] %v: %v", topic, part, off, err)
			} else {
				messages = append(messages, msg)
			}
		}
	}
	if len(messages) < 1 {
		closeFatal("Error: No Messages Received.\n")
	}
	return messages
}

func tailMSGs(flags MSGFlags, topics ...string) []*kafkactl.Message {
	var messages []*kafkactl.Message
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
		endMap := make(map[int32]int64, len(parts))
		offset := getTailValue(flags.Tail)
		for _, p := range parts {
			off, err := client.GetOffsetNewest(topic, p)
			if err != nil {
				closeFatal("Error validating Partition: %v for topic: %v\n", p, err)
			}
			startMap[p] = off + offset
			endMap[p] = off
		}
		msgChan := make(chan *kafkactl.Message, 100)
		doneChan := make(chan bool, len(parts))
		for _, p := range parts {
			go func(topic string, p int32) {
				off := startMap[p]
				for off < endMap[p] {
					msg, err := client.ConsumeOffsetMsg(topic, p, off)
					if err != nil {
						out.Warnf("WARN %v [%v] %v: %v", topic, p, off, err)
					} else {
						msgChan <- msg
					}
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
	}
	if len(messages) < 1 {
		closeFatal("Error: No Messages Received.\n")
	}
	return messages
}

/*
func GetMSG(topic string, partition int32, offset int64) *kafkactl.Message {
	msg, err := client.ConsumeOffsetMsg(topic, partition, offset)
	if err != nil {
		closeFatal("Error retrieving message: %v\n", err)
	}
	return msg
}

func getMSGByTime(topic string, partition int32, datetime string) *kafkactl.Message {
	msg, err := client.OffsetMsgByTime(topic, partition, datetime)
	if err != nil {
		if strings.Contains(err.Error(), "parsing time") && strings.Contains(err.Error(), "cannot parse") {
			errMsg := fmt.Sprintf(`datetime parse error: format should be in the form "mm/dd/YYYY HH:MM:SS.000".`)
			closeFatal("Error retrieving message:\n  %v", errMsg)
		}
		closeFatal("Error retrieving message: %v\n", err)
	}
	return msg
}
*/
