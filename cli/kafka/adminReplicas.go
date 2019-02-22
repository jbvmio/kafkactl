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
	"encoding/json"
	"sort"

	"github.com/jbvmio/kafkactl/cli/zookeeper"

	kafkactl "github.com/jbvmio/kafka"
)

type OpsReplicaFlags struct {
	DryRun            bool
	AllParts          bool
	Brokers           []int32
	Partitions        []int32
	ReplicationFactor int
}

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

type ReplicaDetails struct {
	TopicMetadata []kafkactl.TopicMeta
}

const (
	raPath = `/admin/reassign_partitions`
)

func GetTopicReplicas(topics ...string) ReplicaDetails {
	tMeta := SearchTopicMeta(topics...)
	return ReplicaDetails{TopicMetadata: tMeta}
}

func SetTopicReplicas(flags OpsReplicaFlags, topics ...string) RAPartList {
	var rapList RAPartList
	exact = true
	switch {
	case len(flags.Brokers) < 1 && flags.ReplicationFactor == 0:
		closeFatal("Error: Must Specify either --brokers or --replicas.")
	case len(flags.Brokers) > 0 && flags.ReplicationFactor != 0:
		closeFatal("Error: Cannot use both --brokers and --replicas.")
	case len(flags.Brokers) > 0:
		checkPRE()
		switch {
		case flags.AllParts && len(flags.Partitions) > 0:
			closeFatal("Error: Cannot use both --partitions and --allparts.")
		case flags.AllParts:
			validateBrokers(flags.Brokers)
			rapList = movePartitions(SearchTopicMeta(topics...), flags.Brokers)
		case len(flags.Partitions) > 0:
			validateBrokers(flags.Brokers)
			rapList = movePartitions(validateTopicParts(flags.Partitions, topics...), flags.Brokers)
		default:
			closeFatal("Error: Must Specify either --partitions or --allparts.")
		}
	case flags.ReplicationFactor != 0:
		checkPRE()
		switch {
		case flags.AllParts && len(flags.Partitions) > 0:
			closeFatal("Cannot use both --partitions and --allparts.")
		case flags.AllParts:
			rapList = changeTopicRF(SearchTopicMeta(topics...), flags.ReplicationFactor)
		case len(flags.Partitions) > 0:
			rapList = changeTopicRF(validateTopicParts(flags.Partitions, topics...), flags.ReplicationFactor)
		default:
			closeFatal("Error: Must Specify either --partitions or --allparts.")
		}
	}
	return rapList
}

func ZkCreateRAP(rapList RAPartList) bool {
	data, err := json.Marshal(rapList)
	if err != nil {
		closeFatal("Error on Marshal: %v\n", err)
	}
	handleC("%v", zookeeper.KafkaZK(targetContext, verbose))
	check, err := zookeeper.ZKCheckExists(prePath)
	handleC("Error: %v", err)
	if check {
		closeFatal("Reassign Partitions Already in Progress.")
	}
	zookeeper.ZKCreate(raPath, true, false, data...)
	return true
}

func validateBrokers(brokers []int32) {
	pMap := make(map[int32]bool, len(brokers))
	for _, p := range brokers {
		if pMap[p] {
			closeFatal("Error: invalid broker entered or duplicate.")
		} else {
			pMap[p] = true
		}
	}
	cm, err := client.GetClusterMeta()
	if err != nil {
		closeFatal("Error validating brokers: %v\n", err)
	}
	bMap := make(map[int32]bool, len(brokers))
	for _, b := range brokers {
		for _, broker := range cm.BrokerIDs {
			if broker == b {
				bMap[b] = true
			}
		}
	}
	for _, b := range brokers {
		if !bMap[b] {
			closeFatal("Error validating broker: Invalid Broker > %v\n", b)
		}
	}
}

func validateTopicParts(parts []int32, topics ...string) []kafkactl.TopicMeta {
	var filterTMeta []kafkactl.TopicMeta
	pMap := make(map[int32]bool, len(parts))
	for _, p := range parts {
		if pMap[p] {
			closeFatal("Error: invalid partition entered or duplicate.")
		} else {
			pMap[p] = true
		}
	}
	tMeta := SearchTopicMeta(topics...)
	if len(tMeta) < 1 {
		closeFatal("Error validating topics: %v\n", topics)
	}
	for _, topic := range topics {
		tMap := make(map[int32]bool, len(parts))
		for _, p := range parts {
			for _, tm := range tMeta {
				if tm.Topic == topic && tm.Partition == p {
					tMap[p] = true
					filterTMeta = append(filterTMeta, tm)
				}
			}
		}
		for _, p := range parts {
			if !tMap[p] {
				closeFatal("Error validating topic %v partition: Invalid Partition -> %v", topic, p)
			}
		}
	}

	return filterTMeta
}

func movePartitions(topicMeta []kafkactl.TopicMeta, brokers []int32) RAPartList {
	var raparts []RAPartition
	for _, tm := range topicMeta {
		rap := RAPartition{
			Topic:     tm.Topic,
			Partition: tm.Partition,
			Replicas:  brokers,
		}
		raparts = append(raparts, rap)
	}
	rapList := RAPartList{
		Version:    1,
		Partitions: raparts,
	}
	return rapList
}

func changeTopicRF(tMeta []kafkactl.TopicMeta, rFactor int) RAPartList {
	if len(tMeta) < 1 {
		closeFatal("No results given.")
	}
	tom := client.MakeTopicOffsetMap(tMeta)
	brokers, err := client.GetClusterMeta()
	if err != nil {
		closeFatal("Error retrieving metadata: %v\n", err)
	}
	if rFactor > len(brokers.BrokerIDs) {
		closeFatal("Invalid Number of Brokers Available.\n")
	}

	var raparts []RAPartition
	for _, t := range tom {
		//rap := getRAParts(topic, rFactor, tMeta, brokers)
		raparts = append(raparts, getRAParts(t.Topic, rFactor, t.TopicMeta, brokers)...)
	}
	if len(raparts) < 1 {
		closeFatal("No Changes, Nothing to Assign.")
	}

	rapList := RAPartList{
		Version:    1,
		Partitions: raparts,
	}
	return rapList
}

func getRAParts(topic string, rFactor int, tMeta []kafkactl.TopicMeta, brokers kafkactl.ClusterMeta) []RAPartition {
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
		closeFatal("Invalid Number of Brokers Available.\n")
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
			if len(tm.Replicas) > rFactor {
				for i := 0; i < rFactor; i++ {
					rap.Replicas = append(rap.Replicas, tm.Replicas[i])
				}
				raparts = append(raparts, rap)
			}
		}
	}
	return raparts
}

func sortBrokerByReps(sl []BrokerReplicas) {
	sort.Slice(sl, func(i, j int) bool {
		return sl[i].ReplicaCount < sl[j].ReplicaCount
	})
}

func checkPRE(topics ...string) {
	var preMeta []kafkactl.TopicMeta
	tMeta := SearchTopicMeta(topics...)
	for _, tm := range tMeta {
		if len(tm.Replicas) > 0 {
			if tm.Leader != tm.Replicas[0] {
				preMeta = append(preMeta, tm)
			}
		}
	}
	if len(preMeta) > 1 {
		closeFatal("Error: Preferred Replica Election Needed.")
	}
}
