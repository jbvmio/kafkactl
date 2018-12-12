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
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/spf13/cast"
)

type RAPartList struct {
	Version    int           `json:"version"`
	Partitions []RAPartition `json:"partitions"`
}

type RAPartition struct {
	Topic     string  `json:"topic"`
	Partition int32   `json:"partition"`
	Replicas  []int32 `json:"replicas"`
}

type BrokerReplicas struct {
	BrokerID     int32
	ReplicaCount int32
}

func performPartitionReAssignment(topic string, rFactor int) {
	j := changeTopicRF(topic, rFactor)
	fmt.Printf("%s", j)
	zkCreateReassignPartitions("/admin/reassign_partitions", j)
}

// This is WiP*
func changePartitionCount(topic string, count int32) {
	exact = true
	tMeta := searchTopicMeta(topic)
	if len(tMeta) < 1 {
		log.Fatalf("No results found for topic: %v\n", topic)
	}
	existingPartCount := cast.ToInt32(len(tMeta))
	if count <= existingPartCount {
		log.Fatalf("Invalid Partitions Value: %v\n", count)
	}
	delta := count - existingPartCount
	totalParts := make([][]int32, count)
	partMap := make(map[int32]bool)
	repMap := make(map[int32]bool)
	var lowestReps int
	var replicaChoices []int32
	for _, tm := range tMeta {
		totalParts[tm.Partition] = tm.Replicas
		partMap[tm.Partition] = true
		var rCount int
		for _, r := range tm.Replicas {
			if !repMap[r] {
				replicaChoices = append(replicaChoices, r)
			}
			rCount++
		}
		if lowestReps != 0 {
			lowestReps = rCount
		}
		if lowestReps < rCount {
			lowestReps = rCount
		}
	}
	var i int32
	repUsed := make(map[int32]bool)
	for i < delta {
		var partNum int32
		for x := 0; x < len(totalParts); x++ {
			if !partMap[partNum] {
				var ass []int32
				var success bool
				for _, rep := range replicaChoices {
					if !repUsed[rep] {
						ass = append(ass, rep)
						repUsed[rep] = true
						success = true
						break
					}
				}
				if !success {
					repUsed = nil
					for _, rep := range replicaChoices {
						if !repUsed[rep] {
							ass = append(ass, rep)
							repUsed[rep] = true
							success = true
							break
						}
					}
				}
				totalParts[partNum] = ass
			}
			partNum++
		}
		i++
	}
	err := client.AddPartitions(topic, count, totalParts)
	if err != nil {
		log.Fatalf("Error changing partition count: %v\n", err)
	}
}

func changeTopicRF(topic string, rFactor int) []byte {
	if rFactor < 1 {
		log.Fatalf("Invalid Replication Factor Value: %v\n", rFactor)
	}
	exact = true
	tMeta := searchTopicMeta(topic)
	if len(tMeta) < 1 {
		log.Fatalf("No results found for topic: %v\n", topic)
	}
	brokers, err := client.GetClusterMeta()
	if err != nil {
		log.Fatalf("Error retrieving metadata: %v\n", err)
	}
	if rFactor > len(brokers.BrokerIDs) {
		log.Fatalf("Invalid Number of Brokers Available.\n")
	}

	// Break Out Here Later func(topic, rFactor, tMeta, BR) //
	var BR []BrokerReplicas
	var mostReplicas int32
	for _, b := range brokers.BrokerIDs {
		br := BrokerReplicas{
			BrokerID: b,
		}
		for _, tm := range tMeta {
			for _, r := range tm.Replicas {
				if r == b {
					br.ReplicaCount++
				}
			}
		}
		if br.ReplicaCount > mostReplicas {
			mostReplicas = br.ReplicaCount
		}
		BR = append(BR, br)
	}
	if len(BR) < 2 {
		log.Fatalf("Invalid Number of Brokers Available.\n")
	}
	sortBrokerByReps(BR)
	var raparts []RAPartition
	var previous int32 = -7
	var history []int32
	for t := 0; t < len(tMeta); t++ {
		tm := tMeta[t]
		rap := RAPartition{
			Topic:     tm.Topic,
			Partition: tm.Partition,
		}
		if len(tm.Replicas) != rFactor {
			if len(tm.Replicas) < rFactor {
				reps := tm.Replicas
				delta := rFactor - len(tm.Replicas)
				taken := make(map[int32]bool)
				for i := 0; i < delta; i++ {
					history = append(history, previous)
					var candidates []BrokerReplicas
					for x := 0; x < len(BR); x++ {
						var used bool
						for _, r := range tm.Replicas {
							if r == BR[x].BrokerID {
								used = true
							}
						}
						if !used {
							br := BrokerReplicas{
								BrokerID:     BR[x].BrokerID,
								ReplicaCount: BR[x].ReplicaCount,
							}
							candidates = append(candidates, br)
						}
					}
					sortBrokerByReps(candidates)
					var choices []int32
					for _, c := range candidates {
						choices = append(choices, c.BrokerID)
					}
					var bID int32
					var match bool
					var comparePart int
					if t != len(tMeta)-1 {
						comparePart = t + 1
					}
					maybe := make(map[int32]bool)
					var best int32 = -7777
					var toBeat int32
					for _, next := range tMeta[comparePart].Replicas {
						for _, c := range candidates {
							if !taken[c.BrokerID] {
								if c.BrokerID == next {
									maybe[c.BrokerID] = false
									if best == c.BrokerID {
										best = -7777
									}
								} else {
									maybe[c.BrokerID] = true
									if c.ReplicaCount <= mostReplicas {
										var inHist bool
										if len(history) >= rFactor {
											for _, h := range history[len(history)-rFactor:] {
												if h == c.BrokerID {
													inHist = true
												}
											}
										} else {
											for _, h := range history {
												if h == c.BrokerID {
													inHist = true
												}
											}
										}
										if !inHist {
											if best != -7777 {
												if c.ReplicaCount < toBeat {
													best = c.BrokerID
													toBeat = c.ReplicaCount
												}
											} else {
												best = c.BrokerID
												toBeat = c.ReplicaCount
											}
										}
									}
								}
							}
						}
					}
					if best != -7777 {
						match = true
						bID = best
					} else {
						for k := range maybe {
							if maybe[k] {
								match = true
								if k != previous {
									bID = k
									break
								}
								bID = k
							}
						}
					}
					if !match {
						for _, c := range candidates {
							if !taken[c.BrokerID] {
								if c.BrokerID != previous {
									bID = c.BrokerID
									break
								}
								bID = c.BrokerID
							}
						}
					}
					taken[bID] = true
					reps = append(reps, bID)
					for x := 0; x < len(BR); x++ {
						if BR[x].BrokerID == bID {
							BR[x].ReplicaCount++
							break
						}
					}
					sortBrokerByReps(BR)
					previous = bID
				}
				rap.Replicas = reps
				raparts = append(raparts, rap)
				for _, br := range BR {
					if br.ReplicaCount > mostReplicas {
						mostReplicas = br.ReplicaCount
					}
				}
			}
		}
	}
	if len(raparts) < 1 {
		client.Close()
		log.Fatalf("Nothing to Assign,\n")
	}
	rapList := RAPartList{
		Version:    1,
		Partitions: raparts,
	}
	j, err := json.Marshal(rapList)
	if err != nil {
		log.Fatalf("Error on Marshal: %v\n", err)
	}
	return j
}

func sortBrokerByReps(sl []BrokerReplicas) {
	sort.Slice(sl, func(i, j int) bool {
		return sl[i].ReplicaCount < sl[j].ReplicaCount
	})
}
