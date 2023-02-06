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
	"sort"
	"strings"

	kafkactl "github.com/jbvmio/kafka"
)

type TopicFlags struct {
	FindPRE  bool
	Describe bool
	Group    bool
	Lag      bool
	Leaders  []int32
}

func SearchTopicMeta(topics ...string) []kafkactl.TopicMeta {
	var topicMeta []kafkactl.TopicMeta
	var err error
	switch true {
	case len(topics) < 1:
		topicMeta, err = client.GetTopicMeta()
		if err != nil {
			closeFatal("Error getting topic metadata: %s\n", err)
		}
		if len(topicMeta) < 1 {
			closeFatal("No Topics Seem to Exist.")
		}
	default:
		tMeta, err := client.GetTopicMeta()
		if err != nil {
			closeFatal("Error getting topic metadata: %s\n", err)
		}
		if len(tMeta) < 1 {
			closeFatal("No Topics Seem to Exist.")
		}
		subMatch := true
		switch subMatch {
		case exact:
			for _, t := range topics {
				for _, m := range tMeta {
					if m.Topic == t {
						topicMeta = append(topicMeta, m)
					}
				}
			}
		default:
			for _, t := range topics {
				for _, m := range tMeta {
					if strings.Contains(m.Topic, t) {
						topicMeta = append(topicMeta, m)
					}
				}
			}
		}
	}
	if len(topicMeta) < 1 {
		closeFatal("No Results Found for Topic: %v\n", topics)
	}
	sort.Slice(topicMeta, func(i, j int) bool {
		if topicMeta[i].Topic < topicMeta[j].Topic {
			return true
		}
		if topicMeta[i].Topic > topicMeta[j].Topic {
			return false
		}
		return topicMeta[i].Partition < topicMeta[j].Partition
	})
	return topicMeta
}

func SearchTOM(topics ...string) []kafkactl.TopicOffsetMap {
	tom := GetTopicOffsetMap(SearchTopicMeta(topics...))
	if len(tom) < 1 {
		closeFatal("no results for that group/topic combination\n")
	}
	return tom
}

func GetTopicOffsetMap(tm []kafkactl.TopicMeta) []kafkactl.TopicOffsetMap {
	return client.MakeTopicOffsetMap(tm)
}

func FilterTOMByLeader(tom []kafkactl.TopicOffsetMap, leaders []int32) []kafkactl.TopicOffsetMap {
	validateLeaders(leaders)
	var TOM []kafkactl.TopicOffsetMap
	done := make(map[string]bool)
	for _, t := range tom {
		var topicMeta []kafkactl.TopicMeta
		if !done[t.Topic] {
			done[t.Topic] = true
			for _, tm := range t.TopicMeta {
				for _, leader := range leaders {
					if t.PartitionLeaders[tm.Partition] == leader {
						topicMeta = append(topicMeta, tm)
					}
				}
			}
		}
		tom := kafkactl.TopicOffsetMap{
			Topic:            t.Topic,
			TopicMeta:        topicMeta,
			PartitionOffsets: t.PartitionOffsets,
			PartitionLeaders: t.PartitionLeaders,
		}
		TOM = append(TOM, tom)
	}
	return TOM
}

// FilterTOMByPartitions needs testing ...
func FilterTOMByPartitions(tom []kafkactl.TopicOffsetMap, partitions []int32) []kafkactl.TopicOffsetMap {
	var TOM []kafkactl.TopicOffsetMap
	done := make(map[string]struct{})
	for _, t := range tom {
		var topicMeta []kafkactl.TopicMeta
		partOffsets := make(map[int32]int64)
		leaderOffsets := make(map[int32]int32)
		if _, there := done[t.Topic]; !there {
			done[t.Topic] = struct{}{}
			for _, tm := range t.TopicMeta {
				for _, part := range partitions {
					if _, ok := t.PartitionOffsets[part]; ok {
						topicMeta = append(topicMeta, tm)
						partOffsets[tm.Partition] = t.PartitionOffsets[tm.Partition]
						leaderOffsets[part] = t.PartitionLeaders[part]
					}
				}
			}
		}
		tom := kafkactl.TopicOffsetMap{
			Topic:            t.Topic,
			TopicMeta:        topicMeta,
			PartitionOffsets: partOffsets,
			PartitionLeaders: leaderOffsets,
		}
		TOM = append(TOM, tom)
	}
	return TOM
}

type ValidOffsetTOM struct {
	Tom          []kafkactl.TopicOffsetMap
	ValidOffsets map[string]map[int32]int64
}

func GetValidOffsets(tom []kafkactl.TopicOffsetMap) map[string]map[int32]int64 {
	validOffsets := make(map[string]map[int32]int64)
	for _, topic := range tom {
		for part, offset := range topic.PartitionOffsets {
			vOff, err := GetLastValidOffset(topic.Topic, part, offset)
			if err != nil {
				closeFatal("error retrieving last valid offset: %v\n", err)
			}
			if validOffsets[topic.Topic] == nil {
				validOffsets[topic.Topic] = make(map[int32]int64)
			}
			validOffsets[topic.Topic][part] = vOff
		}
	}
	return validOffsets
}

func validateLeaders(leaders []int32) {
	pMap := make(map[int32]bool, len(leaders))
	for _, p := range leaders {
		if pMap[p] {
			closeFatal("Error: invalid leader/brokerID entered or duplicate.")
		} else {
			pMap[p] = true
		}
	}
}
