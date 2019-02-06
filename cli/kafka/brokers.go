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
)

type Broker struct {
	Address            string
	ID                 int32
	GroupsCoordinating int64
	LeaderReplicas     int64
	PeerReplicas       int64
	TotalReplicas      int64
	MigratingReplicas  int64
	Overload           int64
}

func GetBrokerInfo(b ...string) []*Broker {
	var brokers []*Broker
	gl, err := client.GetGroupListMeta()
	if err != nil {
		if len(gl) > 0 {
			closeFatal("Error getting group metadata: %s\n", err)
		}
	}
	tMeta, err := client.GetTopicMeta()
	if err != nil {
		closeFatal("Error getting topic metadata: %s\n", err)
	}
	brIDMap, err := client.BrokerIDMap()
	if len(b) > 0 {
		tmpMap := make(map[int32]string)
		match := true
		switch match {
		case exact:
			for _, a := range b {
				for id, addr := range brIDMap {
					if a == addr {
						tmpMap[id] = addr
					}
				}
			}
		default:
			for _, a := range b {
				for id, addr := range brIDMap {
					if strings.Contains(addr, a) {
						tmpMap[id] = addr
					}
				}
			}
		}
		brIDMap = tmpMap
	}
	if len(brIDMap) < 1 {
		closeFatal("No Results Found.\n")
	}
	brMap := make(map[int32]*Broker, len(brIDMap))
	if err != nil {
		closeFatal("Error getting broker metadata: %s\n", err)
	}
	for id, addr := range brIDMap {
		br := Broker{
			Address: addr,
			ID:      id,
		}
		brMap[id] = &br
	}
	for _, tm := range tMeta {
		if brMap[tm.Leader] != nil {
			brMap[tm.Leader].LeaderReplicas++
			for _, reps := range tm.Replicas {
				if reps != tm.Leader {
					brMap[tm.Leader].PeerReplicas++
				}
			}
			if len(tm.Replicas) > len(tm.ISRs) {
				brMap[tm.Leader].MigratingReplicas++
			}
			if len(tm.Replicas) > 0 {
				if tm.Leader != tm.Replicas[0] {
					brMap[tm.Leader].Overload++
				}
			}
		}
	}
	for _, grp := range gl {
		if brMap[grp.CoordinatorID] != nil {
			brMap[grp.CoordinatorID].GroupsCoordinating++
		}
	}
	for _, br := range brMap {
		br.TotalReplicas = br.LeaderReplicas + br.PeerReplicas
		brokers = append(brokers, br)
	}
	sort.Slice(brokers, func(i, j int) bool {
		return brokers[i].ID < brokers[j].ID
	})
	return brokers
}
