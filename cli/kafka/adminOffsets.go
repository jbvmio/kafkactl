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

	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/abstraction/kafka"
)

type OpsOffsetFlags struct {
	AllParts       bool
	OffsetNewest   bool
	OffsetOldest   bool
	ShowGroups     bool
	Delete         bool
	Group          string
	Partitions     []int32
	Offset         int64
	RelativeOffset int64
}

type OffsetDetails struct {
	Details        []OffsetDetail
	IncludesGroups bool
}

type OffsetDetail struct {
	Topic             string
	Group             string
	Partition         int32
	TopicOffsetOldest int64
	TopicOffsetNewest int64
	GroupOffset       int64
	Lag               int64
}

type offsetReset struct {
	topic     string
	group     string
	partition int32
	offset    int64
}

func GetTopicOffsets(topics ...string) (topicDetails OffsetDetails) {
	var details []OffsetDetail
	TOM := SearchTOM(topics...)
	for _, t := range TOM {
		var parts []int32
		for k := range t.PartitionOffsets {
			parts = append(parts, k)
		}
		for _, p := range parts {
			d := OffsetDetail{
				Topic:             t.Topic,
				Partition:         p,
				TopicOffsetOldest: t.PartitionOldest[p],
				TopicOffsetNewest: t.PartitionOffsets[p],
			}
			details = append(details, d)
		}
	}
	sort.Slice(details, func(i, j int) bool {
		if details[i].Topic < details[j].Topic {
			return true
		}
		if details[i].Topic > details[j].Topic {
			return false
		}
		return details[i].Partition < details[j].Partition
	})
	topicDetails.Details = details
	topicDetails.IncludesGroups = false
	return
}

func MatchGroupOffsets(topicDetails OffsetDetails) (groupDetails OffsetDetails) {
	var details []OffsetDetail
	groupMetaMap := make(map[string][]kafkactl.GroupMeta)
	groupPartMap := make(map[string][]int32)
	groupDone := make(map[string]bool)
	topicOffsets := topicDetails.Details
	for _, topic := range topicOffsets {
		switch {
		case len(groupMetaMap[topic.Topic]) < 1 || !groupDone[topic.Topic]:
			groupMetaMap[topic.Topic] = GroupMetaByTopics(topic.Topic)
			groupDone[topic.Topic] = true
			fallthrough
		default:
			groupPartMap[topic.Topic] = append(groupPartMap[topic.Topic], topic.Partition)
		}
	}
	for t := range groupDone {
		groupDone[t] = false
	}
	groupOffsetMap := make(map[string]kafkactl.GroupOffsetMap)
	for _, topic := range topicOffsets {
		for _, gmeta := range groupMetaMap[topic.Topic] {
			switch {
			case len(groupOffsetMap[topic.Topic].PartitionOffset) < 1 || !groupDone[topic.Topic]:
				groupOffsetMap[topic.Topic], errd = client.OffSetAdmin().Topic(topic.Topic).Group(gmeta.Group).GetGroupOffsets(groupPartMap[topic.Topic])
				out.IfErrf(errd)
				fallthrough
			default:
				if len(groupOffsetMap[topic.Topic].PartitionOffset) > 0 && !groupDone[topic.Topic] {
					for p := range groupOffsetMap[topic.Topic].PartitionOffset {
						for _, to := range topicOffsets {
							if to.Topic == topic.Topic && to.Partition == p {
								if groupOffsetMap[topic.Topic].PartitionOffset[p] == -1 {
									groupOffsetMap[topic.Topic].PartitionOffset[p] = to.TopicOffsetNewest
								}
								d := OffsetDetail{
									Topic:             groupOffsetMap[topic.Topic].Topic,
									Group:             groupOffsetMap[topic.Topic].Group,
									GroupOffset:       groupOffsetMap[topic.Topic].PartitionOffset[p],
									Partition:         p,
									TopicOffsetOldest: to.TopicOffsetOldest,
									TopicOffsetNewest: to.TopicOffsetNewest,
									Lag:               to.TopicOffsetNewest - groupOffsetMap[topic.Topic].PartitionOffset[p],
								}
								details = append(details, d)
							}
						}
					}
				}
				groupDone[topic.Topic] = true
			}
		}
	}
	sort.Slice(details, func(i, j int) bool {
		if details[i].Topic < details[j].Topic {
			return true
		}
		if details[i].Topic > details[j].Topic {
			return false
		}
		return details[i].Partition < details[j].Partition
	})
	groupDetails.Details = details
	groupDetails.IncludesGroups = true
	return
}

func SetOffsets(flags OpsOffsetFlags, topics ...string) {
	if len(flags.Partitions) < 1 && !flags.AllParts {
		closeFatal("Invalid Number of Partitions entered.")
	}
	for _, topic := range topics {
		switch {
		case flags.Group == "":
			out.Infof("Must Specify a Group if attempting an offset reset or use --delete-to option.")
		default:
			switch {
			case flags.OffsetNewest && flags.OffsetOldest:
				closeFatal("Cannot use both Newest and Oldest.")
			case flags.Offset != 0 && flags.RelativeOffset != 0:
				closeFatal("Cannot specify both offset and reletive offset.")
			case flags.RelativeOffset != 0:
				off := getTailValue(flags.RelativeOffset)
				if off >= 0 {
					closeFatal("Invalid Relative Offset Specified.")
				}
				setTopicOffsets(topic, flags.Group, off, flags.AllParts, flags.Partitions...)
			case flags.OffsetNewest:
				setTopicOffsets(topic, flags.Group, kafkactl.OffsetNewest, flags.AllParts, flags.Partitions...)
			case flags.OffsetOldest:
				setTopicOffsets(topic, flags.Group, kafkactl.OffsetOldest, flags.AllParts, flags.Partitions...)
			default:
				setTopicOffsets(topic, flags.Group, flags.Offset, flags.AllParts, flags.Partitions...)
			}
		}
	}
}

func getTailValue(arg int64) int64 {
	if arg < 0 {
		return arg
	}
	return arg - (arg * 2)
}

func setTopicOffsets(topic, group string, offset int64, allPartitions bool, partitions ...int32) {
	var details []OffsetDetail
	exact = true
	offsetDetails := GetTopicOffsets(topic)
	if len(offsetDetails.Details) < 1 {
		closeFatal("No Results Found for topic / group combination.")
	}
	switch {
	case allPartitions:
		for _, grp := range offsetDetails.Details {
			if grp.Topic == topic {
				grp.Group = group
				details = append(details, grp)
			}
		}
	default:
		for _, part := range partitions {
			for _, grp := range offsetDetails.Details {
				if grp.Topic == topic && grp.Partition == part {
					grp.Group = group
					details = append(details, grp)
				}
			}
		}
	}
	offsetDetails.Details = details
	setOffsets(offsetDetails, offset)
}

func setOffsets(groupDetails OffsetDetails, offset int64) {
	var resetOffsets []offsetReset
	for _, d := range groupDetails.Details {
		or := offsetReset{
			partition: d.Partition,
			group:     d.Group,
			topic:     d.Topic,
		}
		switch {
		case offset == kafkactl.OffsetOldest:
			or.offset = d.TopicOffsetOldest
		case offset == kafkactl.OffsetNewest:
			or.offset = d.TopicOffsetNewest
		case offset < 0:
			or.offset = d.TopicOffsetNewest + offset
		default:
			or.offset = offset
		}
		resetOffsets = append(resetOffsets, or)
	}
	for _, reset := range resetOffsets {
		err := client.OffSetAdmin().Topic(reset.topic).Group(reset.group).ResetOffset(reset.partition, reset.offset)
		out.IfErrf(err)
		if err == nil {
			out.Infof("Successfully sent reset request for group %v topic %v partition %v to offset %v", reset.group, reset.topic, reset.partition, reset.offset)
		}
	}
}

func deleteToOffset(partition int32, offset int64, topics ...string) {
	for _, topic := range topics {
		errd = client.DeleteToOffset(topic, partition, offset)
		handleC("Error deleting to offset: %v", errd)
		out.Infof("Successfully deleted topic %v to offset %v", topic, offset)
	}
}
