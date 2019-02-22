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

	kafkactl "github.com/jbvmio/kafka"
	"github.com/spf13/cast"
)

type LagFlags struct {
	Group bool
	Topic bool
}

type GrpLag struct {
	grpLag kafkactl.GroupLag
	member string
	//memberPartitions map[string][]int32
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

// TotalLag struct def:
type TotalLag struct {
	Group    string
	Topic    string
	TotalLag int64
}

func FindPartitionLag() []PartitionLag {
	grpMeta, err := client.GetGroupMeta()
	handleC("Metadata Error: %v", err)
	var partitionLag []PartitionLag
	grpLag := getGroupLag(grpMeta)
	for _, grp := range grpLag {
		if grp.grpLag.TotalLag > 3 {
			for p := range grp.grpLag.PartitionLag {
				if grp.grpLag.PartitionLag[p] > 3 {
					pl := PartitionLag{
						Group:     grp.grpLag.Group,
						Topic:     grp.grpLag.Topic,
						Partition: p,
						Member:    grp.member,
						Offset:    grp.grpLag.PartitionOffset[p],
						Lag:       grp.grpLag.PartitionLag[p],
					}
					partitionLag = append(partitionLag, pl)
				}
			}
		}
	}
	sortPartitionLag(partitionLag)
	return partitionLag
}

func FindTotalLag() []TotalLag {
	grpMeta, err := client.GetGroupMeta()
	handleC("Metadata Error: %v", err)
	var totalLag []TotalLag
	grpLag := getGroupLag(grpMeta)
	for _, grp := range grpLag {
		if grp.grpLag.TotalLag > 3 {
			tl := TotalLag{
				Group:    grp.grpLag.Group,
				Topic:    grp.grpLag.Topic,
				TotalLag: grp.grpLag.TotalLag,
			}
			totalLag = append(totalLag, tl)
		}
	}
	sortTotalLag(totalLag)
	return totalLag
}

func getGroupLag(grpMeta []kafkactl.GroupMeta) []GrpLag {
	var grpLag []GrpLag
	var wg sync.WaitGroup
	for _, gm := range grpMeta {
		wg.Add(1)
		go func(gm kafkactl.GroupMeta) {
			for _, m := range gm.MemberAssignments {
				for topic, partitions := range m.TopicPartitions {
					groupLag, err := client.OffSetAdmin().Group(gm.Group).Topic(topic).GetTotalLag(partitions)
					handleW("WARN: %v", err)
					switch {
					case groupLag.TotalLag > 0:
						wg.Add(1)
						gLag := GrpLag{
							grpLag: groupLag,
							member: m.ClientID,
						}
						grpLag = append(grpLag, gLag)
						wg.Done()
					}
				}
			}
			wg.Done()
		}(gm)
	}
	wg.Wait()
	return grpLag
}

// Deprecated:
func findPartitionLag() []PartitionLag {
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
		switch {
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

func sortTotalLag(totalLag []TotalLag) {
	sort.Slice(totalLag, func(i, j int) bool {
		if totalLag[i].Group < totalLag[j].Group {
			return true
		}
		if totalLag[i].Group > totalLag[j].Group {
			return false
		}
		return totalLag[i].Topic < totalLag[j].Topic
	})
}
