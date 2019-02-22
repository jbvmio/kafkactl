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

	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/kafka"
)

type GroupFlags struct {
	Describe bool
	Lag      bool
	Topic    bool
}

func SearchGroupListMeta(groups ...string) []kafkactl.GroupListMeta {
	var groupListMeta []kafkactl.GroupListMeta
	var err error
	switch true {
	case len(groups) < 1:
		groupListMeta, err = client.GetGroupListMeta()
		if err != nil {
			closeFatal("Error getting grouplist metadata: %s\n", err)
		}
	default:
		glMeta, err := client.GetGroupListMeta()
		if err != nil {
			closeFatal("Error getting grouplist metadata: %s\n", err)
		}
		subMatch := true
		switch subMatch {
		case exact:
			for _, g := range groups {
				for _, m := range glMeta {
					if m.Group == g {
						groupListMeta = append(groupListMeta, m)
					}
				}
			}
		default:
			for _, g := range groups {
				for _, m := range glMeta {
					if strings.Contains(m.Group, g) {
						groupListMeta = append(groupListMeta, m)
					}

				}
			}
		}
	}
	if len(groupListMeta) < 1 {
		closeFatal("No Results Groups Found.\n")
	}
	sort.Slice(groupListMeta, func(i, j int) bool {
		return groupListMeta[i].Group < groupListMeta[j].Group
	})
	return groupListMeta
}

func SearchGroupMeta(group ...string) []kafkactl.GroupMeta {
	var groupMeta []kafkactl.GroupMeta
	grpMeta, err := client.GetGroupMeta()
	if err != nil {
		closeFatal("Error getting group metadata: %s\n", err)
	}
	switch true {
	case len(group) < 1:
		return grpMeta
	case len(group) > 0:
		subMatch := true
		switch subMatch {
		case exact:
			for _, g := range group {
				for _, gm := range grpMeta {
					if gm.Group == g {
						groupMeta = append(groupMeta, gm)
					}
				}
			}
		default:
			for _, g := range group {
				for _, gm := range grpMeta {
					if strings.Contains(gm.Group, g) {
						groupMeta = append(groupMeta, gm)
					}
				}
			}
		}
	}
	if len(groupMeta) < 1 {
		closeFatal("No Group Results Found.")
	}
	return groupMeta
}

func GetGOM(group string, topics ...string) []kafkactl.GroupOffsetMap {
	var GOM []kafkactl.GroupOffsetMap
	TOM := SearchTOM(topics...)
	if len(TOM) < 1 {
		out.Warnf("no results for topic / group combination %v / %v", topics, group)
	}
	for _, tom := range TOM {
		var parts []int32
		for k := range tom.PartitionOffsets {
			parts = append(parts, k)
		}
		gom, err := client.OffSetAdmin().Topic(tom.Topic).Group(group).GetGroupOffsets(parts)
		switch {
		case err != nil:
			out.Warnf("Error getting %v Offsets: %v", group, err)
		case len(gom.PartitionOffset) < 1:
			out.Warnf("no results for topic / group combination %v / %v", tom.Topic, group)
		default:
			GOM = append(GOM, gom)
		}
	}
	if len(GOM) < 1 {
		closeFatal("no results for topic / group combination %v / %v", topics, group)
	}
	return GOM
}

func GroupMetaByTopics(topics ...string) []kafkactl.GroupMeta {
	var groupMeta []kafkactl.GroupMeta
	grpMeta := SearchGroupMeta()
	matchMap := make(map[string]bool)
	topicsAvail, err := client.ListTopics()
	handleC("Metadata Error: %v", err)
	switch true {
	case exact:
		for _, topic := range topics {
			for _, t := range topicsAvail {
				if t == topic {
					matchMap[t] = true
				}
			}
		}
	default:
		for _, topic := range topics {
			for _, t := range topicsAvail {
				if strings.Contains(t, topic) {
					matchMap[t] = true
				}
			}
		}
	}
	for _, gm := range grpMeta {
		var matchCount uint16
		for _, topic := range gm.Topics {
			if matchMap[topic] {
				matchCount++
			}
		}
		switch true {
		case matchCount < 1:
			continue
		default:
			gMeta := gm
			gMeta.MemberAssignments = nil
			for _, ma := range gm.MemberAssignments {
				mm := kafkactl.MemberMeta{
					ClientHost:      ma.ClientHost,
					ClientID:        ma.ClientID,
					Active:          ma.Active,
					TopicPartitions: make(map[string][]int32),
				}
				for t, p := range ma.TopicPartitions {

					if matchMap[t] {

						mm.TopicPartitions[t] = append(mm.TopicPartitions[t], p...)
						gMeta.MemberAssignments = append(gMeta.MemberAssignments, mm)

					}
				}
			}
			groupMeta = append(groupMeta, gMeta)
		}
	}
	if len(groupMeta) < 1 {
		closeFatal("No Group Results Found.")
	}
	return groupMeta
}

func GroupMetaByMember(members ...string) []kafkactl.GroupMeta {
	var groupMeta []kafkactl.GroupMeta
	grpMeta, err := client.GetGroupMeta()
	handleC("Metadata Error: %v", err)
	for _, gm := range grpMeta {
		var gMeta kafkactl.GroupMeta
		var MA []kafkactl.MemberMeta
		switch true {
		case exact:
			for _, ma := range gm.MemberAssignments {
				for _, member := range members {
					if ma.ClientID == member {
						MA = append(MA, ma)
					}
				}
			}
		default:
			for _, ma := range gm.MemberAssignments {
				for _, member := range members {
					if strings.Contains(ma.ClientID, member) {
						MA = append(MA, ma)
					}
				}
			}
		}
		if len(MA) > 0 {
			gMeta = gm
			gMeta.MemberAssignments = MA
			groupMeta = append(groupMeta, gMeta)

		}
	}
	if len(groupMeta) < 1 {
		closeFatal("No Member Results Found.\n")
	}
	return groupMeta
}
