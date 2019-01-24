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

	"github.com/jbvmio/kafkactl"
)

func SearchGroupListMeta(groups ...string) []kafkactl.GroupListMeta {
	var groupListMeta []kafkactl.GroupListMeta
	var err error
	match := true
	switch match {
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
		for _, g := range groups {
			for _, m := range glMeta {
				if exact {
					if m.Group == g {
						groupListMeta = append(groupListMeta, m)
					}
				} else {
					if strings.Contains(m.Group, g) {
						groupListMeta = append(groupListMeta, m)
					}
				}
			}
		}
	}
	sort.Slice(groupListMeta, func(i, j int) bool {
		return groupListMeta[i].Group < groupListMeta[j].Group
	})
	return groupListMeta
}

func searchGroupMeta(group ...string) []kafkactl.GroupMeta {
	grpMeta, err := client.GetGroupMeta()
	if err != nil {
		closeFatal("Error getting group metadata: %s\n", err)
	}
	if len(group) >= 1 {
		if group[0] != "" {
			var groupMeta []kafkactl.GroupMeta
			for _, g := range group {
				for _, gm := range grpMeta {
					if exact {
						if gm.Group == g {
							groupMeta = append(groupMeta, gm)
						}
					} else {
						if strings.Contains(gm.Group, g) {
							groupMeta = append(groupMeta, gm)
						}
					}
				}
			}
			return groupMeta
		}
	}
	return grpMeta
}

func groupMetaByTopic(topic string, grpMeta []kafkactl.GroupMeta) []kafkactl.GroupMeta {
	var groupMeta []kafkactl.GroupMeta
	for _, gm := range grpMeta {
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
				if exact {
					if t == topic {
						mm.TopicPartitions[t] = append(mm.TopicPartitions[t], p...)
						gMeta.MemberAssignments = append(gMeta.MemberAssignments, mm)
					}
				} else {
					if strings.Contains(t, topic) {
						mm.TopicPartitions[t] = append(mm.TopicPartitions[t], p...)
						gMeta.MemberAssignments = append(gMeta.MemberAssignments, mm)
					}
				}

			}
		}
		groupMeta = append(groupMeta, gMeta)
	}
	return groupMeta
}

func groupMetaByMember(member string, grpMeta []kafkactl.GroupMeta) []kafkactl.GroupMeta {
	var groupMeta []kafkactl.GroupMeta
	for _, gm := range grpMeta {
		gMeta := gm
		gMeta.MemberAssignments = nil
		for _, ma := range gm.MemberAssignments {
			if exact {
				if ma.ClientID == member {
					gMeta.MemberAssignments = append(gMeta.MemberAssignments, ma)
				}
			} else {
				if strings.Contains(ma.ClientID, member) {
					gMeta.MemberAssignments = append(gMeta.MemberAssignments, ma)
				}

			}
		}
		groupMeta = append(groupMeta, gMeta)
	}
	return groupMeta
}
