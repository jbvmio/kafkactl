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
	"regexp"
	"strings"
	"time"

	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/kafka"
	"github.com/spf13/cast"
)

type MSGFlags struct {
	Partitions   []string
	Partition    int32
	Offset       int64
	Tail         int64
	TailTouched  bool
	Follow       bool
	FromTime     string
	ToTime       string
	LastDuration string
}

// OffsetRangeMap contains Topics and a Range of Offsets specified from a beginning and end.
type OffsetRangeMap struct {
	Ranges map[string]map[int32][2]int64
}

// GetMessages returns messages from a kafka topic
func GetMessages(flags MSGFlags, topics ...string) []*kafkactl.Message {
	exact = true
	if flags.TailTouched {
		return tailMSGs(flags, topics...)
	}
	return getMSGs(flags, topics...)
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

// GetMsgOffsets returns offsets for the given topics queried by time.
func GetMsgOffsets(flags MSGFlags, topics ...string) OffsetRangeMap {
	var valid bool
	var beginTime, finishTime int64
	exact = true
	var topicOffsets map[string]map[int32][2]int64
	for _, topic := range topics {
		topicSummary := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
		var parts []int32
		switch {
		case len(topicSummary) != 1:
			closeFatal("Error isolating topic: %v\n", topic)
		case flags.Partition != -1:
			parts = append(parts, flags.Partition)
		case len(flags.Partitions) == 0:
			parts = topicSummary[0].Partitions
		default:
			parts = validateParts(flags.Partitions)
		}
		if topicOffsets == nil {
			topicOffsets = make(map[string]map[int32][2]int64)
		}
		if topicOffsets[topic] == nil {
			topicOffsets[topic] = make(map[int32][2]int64)
		}
		for _, p := range parts {
			topicOffsets[topic][p] = [2]int64{}
		}
	}
	switch {
	case flags.LastDuration != "":
		finishTime = (time.Now().Unix() * 1000)
		dur, err := time.ParseDuration(flags.LastDuration)
		if err != nil {
			closeFatal("Error: Invalid Duration Defined: %v\n", err)
		}
		beginTime = finishTime - cast.ToInt64(dur.Seconds()*1000)
	case flags.ToTime != "":
		if flags.FromTime != "" {
			beginTime = getTimeMillis(flags.FromTime)
		}
		finishTime = getTimeMillis(flags.ToTime)
	case flags.FromTime != "":
		if flags.ToTime == "" {
			finishTime = (time.Now().Unix() * 1000)
		} else {
			finishTime = getTimeMillis(flags.ToTime)
		}
		beginTime = getTimeMillis(flags.FromTime)
	default:
		closeFatal("Error: No Query Parameters Defined.\n")
	}
	if beginTime == -7777 {
		closeFatal("Error Parsing Begin Time: %v: %v\n", flags.FromTime, beginTime)
	}
	if finishTime == -7777 {
		closeFatal("Error Parsing Finish Time: %v: %v\n", flags.ToTime, finishTime)
	}
	for topic, parts := range topicOffsets {
		for p := range parts {
			oRange, v := getBeginFinishOffsets(topic, p, beginTime, finishTime)
			if v {
				valid = true
				topicOffsets[topic][p] = oRange
			} else {
				delete(topicOffsets[topic], p)
			}
		}
	}
	if !valid {
		closeFatal("Error: No Valid Results")
	}
	return OffsetRangeMap{Ranges: topicOffsets}
}

func getBeginFinishOffsets(topic string, partition int32, beginTime, finishTime int64) (offsets [2]int64, valid bool) {
	var err error
	offsets[0], err = client.PartitionOffsetByTime(topic, partition, beginTime)
	if err != nil {
		closeFatal("Error retrieving begin offset: %v\n", err)
	}
	if offsets[0] == -1 {
		out.Warnf("WARN: No valid timestamps found for topic: %v on partition: %v", topic, partition)
		return
	}
	var try int64
	var tries int64
	var finishOffset int64
	for ; finishOffset < offsets[0]; try += 500 {
		tries++
		finishOffset, err = client.PartitionOffsetByTime(topic, partition, finishTime-(try*tries))
		if err != nil {
			out.Warnf("WARN: Error retrieving finish offset for topic: %v on partition: %v", topic, partition)
		}
		if tries%10 == 0 {
			tries *= 10
		}
		if tries%10000 == 0 {
			break
		}
	}
	if finishOffset != -1 {
		valid = true
		offsets[1] = finishOffset
	}

	return
}

func getTimeMillis(dateTime string) int64 {
	timeFormat := "01/02/2006 15:04:05.000"
	targetTime := roundTime(dateTime)
	time, err := time.Parse(timeFormat, targetTime)
	if err != nil {
		return -7777
	}
	return (time.Unix() * 1000)
}

func roundTime(targetTime string) string {
	testTargetTime := strings.Split(targetTime, ".")
	if len(testTargetTime) > 1 {
		addZeros := 3 - (len([]rune(testTargetTime[1])))
		if addZeros > 0 {
			for i := 0; i < addZeros; i++ {
				targetTime += "0"
			}
		}
		return targetTime
	}
	if len(testTargetTime) == 1 {
		targetTime += ".000"
	}
	return targetTime
}

func validateParts(partitions []string) []int32 {
	var tParts []int32
	for _, p := range partitions {
		if match, err := regexp.MatchString(`^[0-9]`, p); !match || err != nil {
			if err != nil {
				closeFatal("Partition Error: %v\n", err)
			}
			closeFatal("Error: invalid partition entered or duplicate.")
		}
		tParts = append(tParts, cast.ToInt32(p))
	}
	pMap := make(map[int32]bool, len(tParts))
	for _, p := range tParts {
		if pMap[p] {
			closeFatal("Error: invalid partition entered or duplicate.")
		} else {
			pMap[p] = true
		}
	}
	return tParts
}
