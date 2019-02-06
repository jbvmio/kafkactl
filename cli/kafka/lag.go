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
	"sync"

	"github.com/jbvmio/kafkactl"
	"github.com/spf13/cast"
)

type LagFlags struct {
	Group bool
	Topic bool
}

// PartitionLag struct def:
type PartitionLag struct {
	Group     string
	Topic     string
	Partition int32
	Member    string
	Offset    int64
	Lag       int64
}

func FindLag() []PartitionLag {
	grpMeta, err := client.GetGroupMeta()
	handleC("Metadata Error: %v", err)
	var wg sync.WaitGroup
	var partitionLag []PartitionLag
	plChan := make(chan PartitionLag, 10000)
	var count uint32
	for _, gm := range grpMeta {
		wg.Add(1)
		go func(gm kafkactl.GroupMeta) {
			for _, m := range gm.MemberAssignments {
				for topic, partitions := range m.TopicPartitions {
					groupLag, err := client.OffSetAdmin().Group(gm.Group).Topic(topic).GetTotalLag(partitions)
					handleW("WARN: %v", err)
					switch true {
					case groupLag.TotalLag > 0:
						wg.Add(1)
						for p := range groupLag.PartitionLag {
							if groupLag.PartitionLag[p] > 0 {
								count++
								pl := PartitionLag{
									Group:     gm.Group,
									Topic:     topic,
									Partition: p,
									Member:    m.ClientID,
									Offset:    groupLag.PartitionOffset[p],
									Lag:       groupLag.PartitionLag[p],
								}
								plChan <- pl
							}
						}
						wg.Done()
					}
				}
			}
			wg.Done()
		}(gm)
	}
	wg.Wait()
	var i uint32
	for ; i < count; i++ {
		pl := <-plChan
		partitionLag = append(partitionLag, pl)
	}
	sortPartitionLag(partitionLag)
	return partitionLag
}

func GetGroupLag(grpMeta []kafkactl.GroupMeta) []PartitionLag {
	var wg sync.WaitGroup
	var partitionLag []PartitionLag
	plChan := make(chan PartitionLag, 10000)
	var count uint32
	for _, gm := range grpMeta {
		wg.Add(1)
		go func(gm kafkactl.GroupMeta) {
			for _, m := range gm.MemberAssignments {
				for topic, partitions := range m.TopicPartitions {
					wg.Add(1)
					count = count + cast.ToUint32(len(partitions))
					wg.Done()
					groupLag, err := client.OffSetAdmin().Group(gm.Group).Topic(topic).GetTotalLag(partitions)
					handleW("WARN: %v", err)
					for _, p := range partitions {
						pl := PartitionLag{
							Group:     gm.Group,
							Topic:     topic,
							Partition: p,
							Member:    m.ClientID,
							Offset:    groupLag.PartitionOffset[p],
							Lag:       groupLag.PartitionLag[p],
						}
						plChan <- pl
					}
				}
			}
			wg.Done()
		}(gm)
	}
	wg.Wait()
	var i uint32
	for ; i < count; i++ {
		pl := <-plChan
		partitionLag = append(partitionLag, pl)
	}
	sortPartitionLag(partitionLag)
	return partitionLag
}

func sortPartitionLag(partitionLag []PartitionLag) {
	sort.Slice(partitionLag, func(i, j int) bool {
		switch true {
		case partitionLag[i].Group == partitionLag[j].Group:
			switch true {
			case partitionLag[i].Topic == partitionLag[j].Topic:
				return partitionLag[i].Partition < partitionLag[j].Partition
			default:
				return partitionLag[i].Topic < partitionLag[j].Topic
			}
		case partitionLag[i].Group < partitionLag[j].Group:
			return true

		case partitionLag[i].Group > partitionLag[j].Group:
			return false
		}
		return partitionLag[i].Partition < partitionLag[j].Partition
	})
}

/*
func getPartitionLag(grpMeta []kafkactl.GroupMeta) []PartitionLag {
	var partitionLag []PartitionLag
	for _, gm := range grpMeta {
		for _, m := range gm.MemberAssignments {
			for topic, partitions := range m.TopicPartitions {
				for _, p := range partitions {
					offset, lag, err := client.OffSetAdmin().Group(gm.Group).Topic(topic).GetOffsetLag(p)
					if err != nil {
						lag = -7777
					}
					pl := PartitionLag{
						Group:     gm.Group,
						Topic:     topic,
						Partition: p,
						Member:    m.ClientID,
						Offset:    offset,
						Lag:       lag,
					}
					partitionLag = append(partitionLag, pl)
				}
			}
		}
	}
	return partitionLag
}

func chanGetPartitionLag(grpMeta []kafkactl.GroupMeta) []PartitionLag {
	var partitionLag []PartitionLag
	for _, gm := range grpMeta {
		for _, m := range gm.MemberAssignments {
			for topic, partitions := range m.TopicPartitions {
				plChan := make(chan PartitionLag, 100)
				for _, p := range partitions {
					go goGetPartitionLag(client, gm.Group, topic, m.ClientID, p, plChan)
				}
				for i := 0; i < len(partitions); i++ {
					select {
					case pl := <-plChan:
						partitionLag = append(partitionLag, pl)
					}
				}
			}
		}
	}
	return partitionLag
}

func goGetPartitionLag(client *kafkactl.KClient, group, topic, clientID string, partition int32, plChan chan PartitionLag) {
	offset, lag, err := client.OffSetAdmin().Group(group).Topic(topic).GetOffsetLag(partition)
	if err != nil {
		lag = -7777
	}
	pl := PartitionLag{
		Group:     group,
		Topic:     topic,
		Partition: partition,
		Member:    clientID,
		Offset:    offset,
		Lag:       lag,
	}
	plChan <- pl
}
*/
