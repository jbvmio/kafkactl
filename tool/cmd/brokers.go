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

package cmd

import (
	"sort"
	"strings"
)

type Broker struct {
	Address           string
	ID                int32
	GroupCoordinating int64
	LeaderPartitions  int64
	ReplicaPartitions int64
	TotalPartitions   int64
	NotLeader         int64
}

func getBrokerInfo(b ...string) []*Broker {
	var brokers []*Broker
	gl, err := client.GetGroupListMeta()
	if err != nil {
		closeFatal("Error getting group metadata: %s\n", err)
	}
	tMeta, err := client.GetTopicMeta()
	if err != nil {
		closeFatal("Error getting topic metadata: %s\n", err)
	}
	brIDMap, err := client.BrokerIDMap()
	if len(b) > 0 {
		tmpMap := make(map[int32]string)
		if exact {
			for _, a := range b {
				for id, addr := range brIDMap {
					if a == addr {
						tmpMap[id] = addr
					}
				}
			}
		} else {
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
			brMap[tm.Leader].LeaderPartitions++
			for _, reps := range tm.Replicas {
				if reps != tm.Leader {
					brMap[tm.Leader].ReplicaPartitions++
				}
			}
			if len(tm.Replicas) > 0 {
				if tm.Leader != tm.Replicas[0] {
					brMap[tm.Leader].NotLeader++
				}
			}
		}
	}
	for _, grp := range gl {
		if brMap[grp.CoordinatorID] != nil {
			brMap[grp.CoordinatorID].GroupCoordinating++
		}
	}
	for _, br := range brMap {
		br.TotalPartitions = br.LeaderPartitions + br.ReplicaPartitions
		brokers = append(brokers, br)
	}
	sort.Slice(brokers, func(i, j int) bool {
		return brokers[i].ID < brokers[j].ID
	})
	return brokers
}
